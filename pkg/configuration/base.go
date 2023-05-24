package configuration

const Backstage string = "backstage"

// ServiceConfig is the configuration for the service.
type ServiceConfig struct {
	API ApiConfig `yaml:"api" json:"api"`
}

// ApiConfig is the configuration for the API
type ApiConfig struct {
	URL   string `json:"url" yaml:"url"`
	Token string `json:"token" yaml:"token"`
}
