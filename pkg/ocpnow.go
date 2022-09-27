package pkg

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
	"text/template"
)

type ClusterInfo struct {
	State             string `yaml:"state"`
	URL               string `yaml:"api_url"`
	Name              string `yaml:"name"`
	CName             string `yaml:"cluster_name"`
	CredId            string `yaml:"credentialId"`
	Id                string `yaml:"cluster_id"`
	PubSvcEndpointURL string `yaml:"public_service_endpoint_url"`
	Infra             string `yaml:"infra_host"`
}

type MetaInfo struct {
	Name string `yaml:"project_name"`
}

type CredInfo struct {
	Name   string `yaml:"name"`
	Infra  string `yaml:"infra"`
	State  string `yaml:"state"`
	ApiKey string `yaml:"api_key"`
}

type Project struct {
	Clusters    map[string]ClusterInfo `yaml:"clusters"`
	Meta        MetaInfo               `yaml:"general"`
	Credentials map[string]CredInfo    `yaml:"credentials"`
}

func (p *Project) Write(out io.Writer) error {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	consoleTemplate := `Project "{{.Meta.Name}}"
Clusters:
{{range.Clusters}}
- Name: {{.Name}}{{.CName}} ({{.State}} to {{.Infra | ToUpper}})
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
