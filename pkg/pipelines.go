package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"

	getter "github.com/hashicorp/go-getter"
	logger "github.com/sirupsen/logrus"
)

const RawGitHubUrlHost = "raw.githubusercontent.com"

type ResolverOption uint32

var slugRegex *regexp.Regexp

func (o ResolverOption) Includes(opt ResolverOption) bool {
	return o&opt != 0
}

const (
	UsePipelineDefaults ResolverOption = 1 << iota
	UseEnvironmentVars
	UseCommandLineArgs
	UsePromptAnswers
)

const DefaultParseOptions = UseEnvironmentVars

func IsPipeline(p v1beta1.Pipeline) bool {
	return p.Kind == "Pipeline"
}

func IsPipelineRun(p v1beta1.PipelineRun) bool {
	return p.Kind == "PipelineRun"
}

// BuildDestination builds the destination path using the git path to the
// pipeline.
func BuildDestination(base string, gitURL string) (string, error) {
	if len(base) == 0 {
		return "", fmt.Errorf("base destination is not set")
	}
	u, err := url.Parse(gitURL)
	if err != nil {
		return "", nil
	}
	pipelinePath := u.Path
	return path.Join(base, pipelinePath), nil
}

// MapGitUrlToRaw updates the URL that you see in the browser to one that
// contains the raw git URL. For example, https://github.com/cloud-native-toolkit/deployer-cloud-pak-deployer/blob/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml
// is really https://raw.githubusercontent.com/cloud-native-toolkit/deployer-cloud-pak-deployer/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml
func MapGitUrlToRaw(id string) (string, error) {
	raw, err := url.Parse(id)
	if err != nil {
		return "", err
	}

	if raw.Host == "raw.githubusercontent.com" {
		// Oh, cool, they already gave us the raw URL. Nothing to see here...
		return raw.String(), nil
	}

	raw.Host = "raw.githubusercontent.com"
	pathParts := strings.Split(raw.Path, "/")

	if len(pathParts) < 5 {
		return "", fmt.Errorf("invalid path: %s", raw.Path)
	}

	newParts := make([]string, 0)
	// the org
	newParts = append(newParts, pathParts[1])
	// the repo
	newParts = append(newParts, pathParts[2])
	// skipping three, because that's blob...
	// the branch
	newParts = append(newParts, pathParts[4])
	// the rest
	newParts = append(newParts, pathParts[5:]...)
	raw.Path = strings.Join(newParts, "/")
	return raw.String(), nil
}

// GitServiceClient can download (get) the objects from a Git repository location
type GitServiceClient struct {
	BaseDest string
}

type MarshallerFunc func(b []byte) (interface{}, error)

// Get the Pipeline from the Git repository
func (g *GitServiceClient) Get(gitURL string, marshaller MarshallerFunc) (interface{}, error) {
	dest, err := BuildDestination(g.BaseDest, gitURL)

	if err != nil {
		return nil, err
	}

	err = getter.GetFile(dest, gitURL)

	if err != nil {
		return nil, err
	}

	yamlFile, err := os.ReadFile(dest)
	if err != nil {
		return nil, err
	}

	// if there is no error, unmarshal the pipeline YAML from the file
	return marshaller(yamlFile)
}

func UnmarshalPipeline(yamlFile []byte) (interface{}, error) {
	//yamlDecoder := yaml.NewDecoder(r)
	var pipeline v1beta1.Pipeline
	found := false
	for {
		//err := yamlDecoder.Decode(&pipeline)
		_, _, err := scheme.Codecs.UniversalDeserializer().Decode(yamlFile, nil, &pipeline)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return nil, err
		}
		if IsPipeline(pipeline) {
			found = true
			break
		} else {
			logger.Tracef("Found document of type %s; skipping...", pipeline.Kind)
		}
	}

	if !found {
		return nil, fmt.Errorf("could not find pipeline in file: %s", yamlFile)
	}
	return &pipeline, nil
}

func UnmarshalPipelineRun(content []byte) (interface{}, error) {
	//yamlDecoder := yaml.NewDecoder(r)
	var pr v1beta1.PipelineRun
	found := false
	for {
		//err := yamlDecoder.Decode(&pipeline)
		_, _, err := scheme.Codecs.UniversalDeserializer().Decode(content, nil, &pr)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return nil, err
		}
		if IsPipelineRun(pr) {
			found = true
			break
		} else {
			logger.Tracef("Found document of type %s; skipping...", pr.Kind)
		}
	}

	if !found {
		return nil, fmt.Errorf("could not find pipeline run in file: %s", content)
	}
	return &pr, nil
}

func MergePipelineRun(run *v1beta1.PipelineRun, pl *v1beta1.Pipeline, reader ParamReader, resolver ParamResolver) (*v1beta1.PipelineRun, error) {
	// first, look up the param keys from the reader...
	result := run.DeepCopy()
	parms, err := reader.Params()
	if err != nil {
		return nil, err
	}
	var updated v1beta1.Params = make([]v1beta1.Param, 0)
	//run.Spec.Params.DeepCopyInto(&updated)
	for _, p := range parms {
		// Lookup the param, and if it exists, update it...
		// if it does not exist, add it.
		param, exists := FindParam(updated, p.Name)
		val, found := resolver.Lookup(p.Name)
		if !found {
			logger.Debugf("value of parameter %s was not found", p.Name)
		}
		if exists {
			var update v1beta1.Param
			param.DeepCopyInto(&update)
			update.Value.StringVal = val
			updated = append(updated, update)
		} else {
			updated = append(updated, v1beta1.Param{
				Name: p.Name,
				Value: v1beta1.ParamValue{
					Type:      v1beta1.ParamTypeString,
					StringVal: val,
				},
			})
		}
	}
	// Now iterate through the parameters in the original run and add them along with their values
	// if they aren't in the set
	for _, op := range run.Spec.Params {
		_, found := FindParam(updated, op.Name)
		if !found {
			var update v1beta1.Param
			op.DeepCopyInto(&update)
			updated = append(updated, update)
		}
	}
	result.Spec.Params = updated
	return result, nil
}

func FindParam(in v1beta1.Params, name string) (*v1beta1.Param, bool) {
	for _, p := range in {
		if p.Name == name {
			return &p, true
		}
	}
	return nil, false
}

// ExecPipelineRun
func ExecPipelineRun(pipeline *v1beta1.Pipeline, run *v1beta1.PipelineRun, runScript string, useContainer bool, cluster ClusterInfo, cred CredInfo, in io.Reader, out io.Writer) error {
	// Now serialize the pipeline and the pipeline runs to files
	pipelineURL := HomeTempFile(MustITZHomeDir(), "pipeline.json")
	pipelineRunURL := HomeTempFile(MustITZHomeDir(), "pipelinerun.json")
	if err := WriteToFile(pipelineURL, true, pipeline); err != nil {
		return err
	}
	if err := WriteToFile(pipelineRunURL, true, run); err != nil {
		return err
	}

	err := WriteFile(filepath.Join(MustITZHomeDir(), "cache", "run.sh"), []byte(runScript))
	if err != nil {
		return err
	}
	if useContainer {
		logger.Debugf("Using container to execute the pipeline...")
		return fmt.Errorf("not currently implemented, make sure -c or --use-container is set to false")
	} else {
		logger.Debugf("Using local commands to execute the pipeline...")
		cmd := exec.Command("bash", filepath.Join(MustITZHomeDir(), "cache", "run.sh"))
		cmd.Env = append(cmd.Env, fmt.Sprintf("ITZ_OC_USER=%s", cred.Name))
		cmd.Env = append(cmd.Env, fmt.Sprintf("ITZ_OC_PASS=%s", cred.ApiKey))
		cmd.Env = append(cmd.Env, fmt.Sprintf("ITZ_OC_URL=%s", cluster.URL))
		cmd.Env = append(cmd.Env, fmt.Sprintf("ITZ_PIPELINE=%s", pipelineURL))
		cmd.Env = append(cmd.Env, fmt.Sprintf("ITZ_PIPELINE_RUN=%s", pipelineRunURL))
		logger.Tracef("running command: %s", cmd)
		cmd.Stdout = out
		cmd.Stderr = out
		cmd.Stdin = in
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

type PipelineParamParts struct {
	Description  string
	ParamOptions PipelineParamOptions
}

func (p *PipelineParamParts) HasOptions() bool {
	return len(p.ParamOptions.Options) > 0
}

type PipelineParamOptions struct {
	Options []PipelineParamOption
}

type PipelineParamOption struct {
	Text    string
	Value   string
	Default string
}

func (o *PipelineParamOption) IsDefault() bool {
	val, err := strconv.ParseBool(o.Default)
	if err != nil {
		return false
	}
	return val
}

func ParseParamDescription(from string) (*PipelineParamParts, error) {
	// See https://github.ibm.com/skol/backstage-catalog/blob/main/MODEL.md#getting-parameters-from-the-tekton-pipelines-for-use-in-gui-applications
	// for the format of this field. It should look like this:
	// specify the preferred storageclass
	// {
	//	"options": [
	//    {"text": "thin","value": "thin", "default": "true"}
	//    {"text": "gp2","value": "gp2" }
	//    {"text": "ocs-storagecluster-cephfs","value": "ocs-storagecluster-cephfs" }
	//  ]
	//}
	lines := strings.Split(from, "\n")
	descr := lines[0]
	var options PipelineParamOptions
	if len(lines) > 1 {
		// try to read the rest into the options using the JSON reader
		r := strings.NewReader(strings.Join(lines[1:], "\n"))
		if r.Size() > 0 {
			err := json.NewDecoder(r).Decode(&options)
			if err != nil {
				return &PipelineParamParts{
					Description: descr,
				}, err
			}
		}
	}

	return &PipelineParamParts{
		Description:  descr,
		ParamOptions: options,
	}, nil
}

// ParamResolver resolves the parameter values from a source.
type ParamResolver interface {
	Supports() ResolverOption
	// EnabledFor returns true if the `ParamResolver` is enabled for the given
	// option.
	EnabledFor(opt ResolverOption) bool
	// Lookup returns the value of the parameter, if found, as well as
	// a bool that indicates if it was found.
	Lookup(p string) (string, bool)
}

// ParamReader reads the `ParamSpec` objects from a source.
type ParamReader interface {
	// Params gets all `ParamSpec` objects
	Params() ([]v1beta1.ParamSpec, error)
}

const DefaultPrefix = "ITZ_"

type EnvParamResolver struct {
	Prefix string
}

func (p *EnvParamResolver) Supports() ResolverOption {
	return UseEnvironmentVars
}

func (p *EnvParamResolver) EnabledFor(opt ResolverOption) bool {
	return opt.Includes(p.Supports())
}

func (p *EnvParamResolver) Lookup(k string) (string, bool) {
	return os.LookupEnv(ToEnvVar(p.Prefix, k))
}

func NewEnvParamResolver() ParamResolver {
	return &EnvParamResolver{
		Prefix: DefaultPrefix,
	}
}

type ArgsParamParser struct {
	args   []string
	params map[string]string
}

func (p *ArgsParamParser) Supports() ResolverOption {
	return UseCommandLineArgs
}

func (p *ArgsParamParser) EnabledFor(opt ResolverOption) bool {
	return opt.Includes(p.Supports())
}

func (p *ArgsParamParser) Lookup(k string) (string, bool) {
	val, exists := p.params[k]
	return val, exists
}

func NewArgsParamParser(args []string) ParamResolver {
	paramMap := make(map[string]string, len(args))
	for _, a := range args {
		k := strings.Split(a, "=")
		if len(k) >= 2 {
			paramMap[k[0]] = strings.Join(k[1:], "=")
		} else {
			logger.Tracef("ignoring malformed argument: \"%s\"", a)
		}
	}
	return &ArgsParamParser{
		args:   args,
		params: paramMap,
	}
}

type PipelineResolver struct {
	pipeline *v1beta1.Pipeline
	params   map[string]v1beta1.ParamSpec
}

func (p *PipelineResolver) Supports() ResolverOption {
	return UsePipelineDefaults
}

// EnabledFor returns true if the `ParamResolver` is enabled for the given
// option.
func (p *PipelineResolver) EnabledFor(opt ResolverOption) bool {
	return opt.Includes(p.Supports())
}

// Lookup returns the value of the parameter, if found, as well as
// a bool that indicates if it was found. In the case of the
// `PipelineResolver`, this returns true if the `Pipeline` parameter
// has a default value.
func (p *PipelineResolver) Lookup(k string) (string, bool) {
	param, found := p.params[k]
	if found {
		return param.Default.StringVal, true
	}
	return "", false
}

func (p *PipelineResolver) Params() ([]v1beta1.ParamSpec, error) {
	return p.pipeline.Spec.Params, nil
}

// NewPipelineResolver creates a new pipeline resolver that will resolve the
// parameters so long as they have default values. This is useful when
// accepting the defaults and not prompting the user unnecessarily.
func NewPipelineResolver(p *v1beta1.Pipeline) *PipelineResolver {
	paramMap := make(map[string]v1beta1.ParamSpec)
	for _, param := range p.Spec.Params {
		if param.Default != nil && len(param.Default.StringVal) > 0 {
			paramMap[param.Name] = param
		}
	}
	return &PipelineResolver{
		pipeline: p,
		params:   paramMap,
	}
}

// ChainedResolver is a Resolver itself that resolves the variables from different locations,
// such as a Pipeline defaults, the command line, or environment variables.
type ChainedResolver struct {
	options   ResolverOption
	resolvers []ParamResolver
}

func (p *ChainedResolver) Supports() ResolverOption {
	var opt ResolverOption
	for _, r := range p.resolvers {
		opt = opt | r.Supports()
	}
	return opt
}

// EnabledFor returns true if the `ParamResolver` is enabled for the given
// option.
func (p *ChainedResolver) EnabledFor(opt ResolverOption) bool {
	for _, r := range p.resolvers {
		if r.EnabledFor(opt) {
			return true
		}
	}
	return false
}

func (p *ChainedResolver) Lookup(k string) (string, bool) {
	// Loop through each of the resolvers and return with the correct
	// one given the options.
	for _, r := range p.resolvers {
		if r.EnabledFor(p.options) {
			val, exists := r.Lookup(k)
			logger.Tracef("Using value %s found in resolver for %v", val, r.Supports())
			if exists {
				return val, exists
			}
		}
	}
	return "", false
}

func NewChainedResolver(opt ResolverOption, enabled ...ParamResolver) *ChainedResolver {
	ordered := make([]ParamResolver, len(enabled))
	copy(ordered, enabled)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Supports() > ordered[j].Supports()
	})
	return &ChainedResolver{
		resolvers: ordered,
		options:   opt,
	}
}

type PromptResolver struct {
	prompt *prompt.Prompt
}

func (p *PromptResolver) Supports() ResolverOption {
	return UsePromptAnswers
}

func (p *PromptResolver) EnabledFor(opt ResolverOption) bool {
	return opt.Includes(p.Supports())
}

func (p *PromptResolver) Lookup(k string) (string, bool) {
	param, found := p.prompt.LookupAnswer(k)
	if found {
		return param, true
	}
	return "", false
}

func NewPromptResolver(p *prompt.Prompt) ParamResolver {
	return &PromptResolver{
		prompt: p,
	}
}

func BuildPipelinePrompt(name string, reader ParamReader, resolver ParamResolver) (*prompt.Prompt, error) {
	root, err := prompt.NewPromptBuilder().
		Path("root").
		Text(fmt.Sprintf("Do you want to install %s?", name)).
		WithOptions(prompt.YesNo()).
		Build()
	// Loop through the parameters for the pipeline and add them to the prompt
	params, err := reader.Params()
	if err != nil {
		return nil, err
	}
	for _, param := range params {
		// If the parameter can already be looked up in the resolver, we don't need
		// to bother the user with it. We will trust that the resolvers are the right
		// resolvers.
		if _, exists := resolver.Lookup(param.Name); exists {
			continue
		}
		builder := prompt.NewPromptBuilder().
			Path(Sluggify(param.Name))

		// There is some logic embedded into the pipeline file. We parse the description
		// and, if there are options embedded in the description, we add those options here.
		p, err := ParseParamDescription(param.Description)
		if err != nil {
			// TODO: Do we need more robust error handling here?
			return nil, err
		}

		if param.Default != nil && len(param.Default.StringVal) > 0 {
			builder.WithDefault(param.Default.StringVal)
		}

		builder.Text(p.Description)

		if p.HasOptions() {
			for _, opt := range p.ParamOptions.Options {
				if opt.IsDefault() {
					builder.AddDefaultOption(opt.Text, opt.Value)
				} else {
					builder.AddOption(opt.Text, opt.Value)
				}
			}
			builder.WithValidator(prompt.CaseInsensitveTextOptionValidator)
		}

		q, err := builder.Build()

		if err != nil {
			// TODO: log this better and perhaps do something about it if the build options support it.
			continue
		}
		root.AddSubPrompt(q)
	}
	return root, err
}

func Sluggify(s string) string {
	formatted := slugRegex.ReplaceAllString(s, "-")
	return strings.ToLower(formatted)
}

func ToEnvVar(prefix, k string) string {
	formatted := slugRegex.ReplaceAllString(k, "_")
	return fmt.Sprintf("%s%s", prefix, strings.ToUpper(formatted))
}

// HomeTempFile returns the path name of a file based on the metadata
func HomeTempFile(base string, name string) string {
	// TODO: add a generated directory name
	return filepath.Join(base, "cache", name)
}

func WriteToFile(fn string, create bool, obj runtime.Object) error {
	if create {
		dir := filepath.Dir(fn)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	err := scheme.Codecs.LegacyCodec(v1beta1.SchemeGroupVersion).Encode(obj, buf)
	if err != nil {
		return err
	}
	return WriteFile(fn, buf.Bytes())
}

func init() {
	slugRegex = regexp.MustCompile("[^0-9a-zA-Z]")
}
