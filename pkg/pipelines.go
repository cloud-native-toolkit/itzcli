package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

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
	UseEnvironmentVars ResolverOption = 1 << iota
	UsePipelineDefaults
	UseCommandLineArgs
)

const DefaultParseOptions = UseEnvironmentVars

func IsPipeline(p v1beta1.Pipeline) bool {
	return p.Kind == "Pipeline"
}

type PipelineServiceClient interface {
	Get(id string) (*v1beta1.Pipeline, error)
	GetAll() ([]*v1beta1.Pipeline, error)
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

// Get the Pipeline from the Git repository
func (g *GitServiceClient) Get(gitURL string) (*v1beta1.Pipeline, error) {
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
	pipeline, err := unmarshalPipeline(yamlFile)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func unmarshalPipeline(yamlFile []byte) (*v1beta1.Pipeline, error) {
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

func (g *GitServiceClient) GetAll() ([]*v1beta1.Pipeline, error) {
	panic("not implemented")
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

func (p *EnvParamResolver) EnabledFor(opt ResolverOption) bool {
	return opt.Includes(UseEnvironmentVars)
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

func (p *ArgsParamParser) EnabledFor(opt ResolverOption) bool {
	return opt.Includes(UseCommandLineArgs)
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

// EnabledFor returns true if the `ParamResolver` is enabled for the given
// option.
func (p *PipelineResolver) EnabledFor(opt ResolverOption) bool {
	return opt.Includes(UsePipelineDefaults)
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
			if exists {
				return val, exists
			}
		}
	}
	return "", false
}

func NewChainedResolver(opt ResolverOption, enabled ...ParamResolver) ParamResolver {
	return &ChainedResolver{
		resolvers: enabled,
		options:   opt,
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

func init() {
	slugRegex = regexp.MustCompile("[^0-9a-zA-Z]")
}
