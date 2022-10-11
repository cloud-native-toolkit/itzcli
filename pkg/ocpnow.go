package pkg

import (
	logger "github.com/sirupsen/logrus"
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
	Name  string `json:"name"`
	Value string `json:"default"`
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
	if err == nil {
		tmpl.Execute(out, p)
	}
	return err
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
	// and, for each field, look for a tag (default "tfvar" and, if found, we will
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
