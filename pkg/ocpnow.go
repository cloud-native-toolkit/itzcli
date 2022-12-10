package pkg

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"reflect"
	"strings"
	"text/template"
)

type ClusterInfo struct {
	State             string `yaml:"state"`
	URL               string `yaml:"api_url" tfvar:"server_url"`
	Name              string `yaml:"name"`
	CName             string `yaml:"cluster_name"`
	CredId            string `yaml:"credentialId"`
	Id                string `yaml:"cluster_id"`
	PubSvcEndpointURL string `yaml:"public_service_endpoint_url"`
	Infra             string `yaml:"infra_host"`
	Region            string `yaml:"region" tfvar:"region"`
}

type MetaInfo struct {
	Name string `yaml:"project_name"`
}

type CredInfo struct {
	Name   string `yaml:"name"`
	Infra  string `yaml:"infra"`
	State  string `yaml:"state"`
	ApiKey string `yaml:"api_key" tfvar:"ibmcloud_api_key"`
}

type Project struct {
	Clusters    map[string]ClusterInfo `yaml:"clusters"`
	Meta        MetaInfo               `yaml:"general"`
	Credentials map[string]CredInfo    `yaml:"credentials"`
}

type JobParam struct {
	Name    string `json:"name"`
	Value   string `json:"value,omitempty"`
	Default string `json:"default,omitempty"`
}

func Lookup(parm JobParam, vars map[string]string) (string, bool) {
	for k, v := range vars {
		if strings.EqualFold(parm.Name, k) {
			if len(v) > 0 {
				return v, true
			}
		}
	}
	return "", false
}

func (p *Project) Write(out io.Writer) error {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	consoleTemplate := `Project "{{.Meta.Name}}"
Clusters:
{{range.Clusters}}
- Name: {{.Name}}{{.CName}} ({{.State}} to {{.Infra | ToUpper}})
  Id: {{.Id}}
  URL: {{.URL}}
{{end}}
`
	tmpl, err := template.New("Project").Funcs(funcMap).Parse(consoleTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(out, p)
}

func FindClusterByName(in *Project, name string) (*string, error) {
	for k, cluster := range in.Clusters {
		if cluster.Name == name {
			return &k, nil
		}
		if cluster.CName == name {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("cluster with name <%s> not found", name)
}

// LoadProject loads the given project from the yaml file.
func LoadProject(path string) (*Project, error) {
	var proj Project
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &proj)
	return &proj, err
}

// ResolveVarConfig are the configuration options for the ResolveVars func.
type ResolveVarConfig struct {
	Prefix string
	Tag    string
}

// ResolveVars looks through the given structure and returns a map of the values
// If cfg is nil, a default set of configuration options (see NewDefaultResolveVarConfig)
// are used.
func ResolveVars(ref interface{}, cfg *ResolveVarConfig) (map[string]string, error) {
	config := NewDefaultResolveVarConfig()
	if cfg != nil {
		config = mergeConfig(config, cfg)
	}
	logger.Tracef("Using config: %v", config)
	vars := make(map[string]string)

	// Borrowing heavily from https://github.com/caarlos0/env/blob/main/env.go for
	// the parsing code here. We're going to go through the struct that we have
	// and, for each field, look for a tag (default "tfvar") and, if found, we will
	// add both the tag's name (with the prefix) to the map coming back along with
	// the string value representation.
	ptr := reflect.ValueOf(ref)
	obj := ptr.Elem()
	t := obj.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		valF := obj.Field(i)
		name, ok := f.Tag.Lookup("tfvar")
		if ok && len(name) > 0 {
			valI := valF.Interface()
			thisVal := reflect.ValueOf(valI)
			if thisVal.Kind() == reflect.String {
				vars[config.Prefix+name] = thisVal.String()
			}
		}
	}

	return vars, nil
}

func mergeConfig(config *ResolveVarConfig, cfg *ResolveVarConfig) *ResolveVarConfig {
	// merge the two configurations...
	result := config
	if cfg.Prefix != "" {
		result.Prefix = cfg.Prefix
	}
	if cfg.Tag != "" {
		result.Tag = cfg.Tag
	}
	return result
}

func NewDefaultResolveVarConfig() *ResolveVarConfig {
	return &ResolveVarConfig{
		Prefix: "TF_VAR_",
		Tag:    "tfvar",
	}
}

// BuildParamResolver is used to resolve missing parameters
type BuildParamResolver struct {
	// The project
	project *Project
	// The name of the cluster from the project.yaml file
	cluster string
	// The job parameters as required by whatever job we're going to run.
	jobParams []JobParam
	// The map of the cluster variables that were found in the ocpnow project
	// config
	clusterVars map[string]string
	// The map of variables that come from the credentials section of the ocpnow
	// project.yaml
	credVars map[string]string
	// The array of variables that we still need to ask the user
	askVars []string
	// The root prompt
	rootPrompt *prompt.Prompt
}

func (r *BuildParamResolver) BuildPrompter(solution string) (*prompt.Prompt, error) {
	// Okay, so now I have the required vars and I can now build up the prompts
	// to ask my user for the values.
	builder := prompt.NewPromptBuilder()

	rootQuestion, err := builder.Path("proceed").
		Text(fmt.Sprintf("This will deploy the solution %s to cluster %s; continue?", solution, r.cluster)).
		WithOptions(prompt.YesNo()).
		Build()

	if err != nil {
		return nil, err
	}

	for _, v := range r.askVars {
		logger.Tracef("Building prompt for <%s>", v)
		subP, _ := prompt.NewPromptBuilder().
			Path(v).
			Textf("What value would you like to use for '%s'?", v).
			Build()
		rootQuestion.AddSubPrompt(subP)
	}
	r.rootPrompt = rootQuestion
	return rootQuestion, nil
}

func (r *BuildParamResolver) ResolvedParams() map[string]string {
	m := make(map[string]string)
	for _, p := range r.jobParams {
		envVal, ok := os.LookupEnv(p.Name)
		if ok {
			logger.Tracef("Using build parameter <%s> from environment with value <%s>.", p.Name, envVal)
			m[p.Name] = fmt.Sprintf("%v", envVal)
		}
	}
	if r.rootPrompt != nil {
		for k, v := range r.rootPrompt.VarMap() {
			logger.Tracef("Adding build parameter <%s> with value <%s>.", k, v)
			m[k] = v
		}
	}
	for k, v := range r.clusterVars {
		logger.Tracef("Adding build parameter <%s> with value <%s>.", k, v)
		m[k] = v
	}
	for k, v := range r.credVars {
		logger.Tracef("Adding build parameter <%s> with value <%s>.", k, v)
		m[k] = v
	}
	return m
}

func NewBuildParamResolver(project *Project, cluster string, params []JobParam) (*BuildParamResolver, error) {

	cRef, err := FindClusterByName(project, cluster)
	if err != nil {
		return nil, err
	}

	cInfo := project.Clusters[*cRef]
	clusterVars, _ := ResolveVars(&cInfo, nil)
	logger.Debugf("Got cluster vars: %v", clusterVars)
	logger.Debugf("Using region: %s", clusterVars["TF_VAR_region"])
	credInfo := project.Credentials[cInfo.CredId]
	credVars, _ := ResolveVars(&credInfo, nil)
	logger.Debugf("Got cred vars: %v", credVars)

	// Now we have a list of the required parameters (vars), and we need
	// to look at the ones that we have and that have values (clusterVars and
	// credVars), and also look into the os.Environment. We'll build a list
	// of the required ones that we don't have values for so that we can
	// prompt the user.

	askVars := make([]string, 0)
	for _, p := range params {
		_, foundInUser := Lookup(p, clusterVars)
		_, foundInCred := Lookup(p, credVars)
		_, foundInEnv := os.LookupEnv(p.Name)
		if !foundInUser && !foundInCred && !foundInEnv && len(p.Value) == 0 {
			logger.Debugf("Found no existance of <%s>, adding to list of required vars", p.Name)
			askVars = append(askVars, p.Name)
		}
	}

	return &BuildParamResolver{
		cluster:     cluster,
		jobParams:   params,
		clusterVars: clusterVars,
		credVars:    credVars,
		askVars:     askVars,
	}, nil
}
