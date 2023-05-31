package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"strings"

	getter "github.com/hashicorp/go-getter"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const RawGitHubUrlHost = "raw.githubusercontent.com"

//region Pipeline

// ObjectMetadata is the metadata of an object
type ObjectMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type PipelineSpec struct {
}

type Pipeline struct {
	Kind     string         `yaml:"kind"`
	Metadata ObjectMetadata `yaml:"metadata"`
	Spec     PipelineSpec   `yaml:"spec"`
}

func (p *Pipeline) Name() string {
	return p.Metadata.Name
}

func (p *Pipeline) IsPipeline() bool {
	return p.Kind == "Pipeline"
}

//endregion Pipeline

type PipelineServiceClient interface {
	Get(id string) (*Pipeline, error)
	GetAll() ([]*Pipeline, error)
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
func (g *GitServiceClient) Get(gitURL string) (*Pipeline, error) {
	pipelineRepo, err := MapGitUrlToRaw(gitURL)
	if err != nil {
		return nil, err
	}

	dest, err := BuildDestination(g.BaseDest, pipelineRepo)

	if err != nil {
		return nil, err
	}

	err = getter.GetFile(dest, pipelineRepo)

	if err != nil {
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(dest)
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

func unmarshalPipeline(yamlFile []byte) (*Pipeline, error) {
	r := bytes.NewReader(yamlFile)
	yamlDecoder := yaml.NewDecoder(r)
	var pipeline Pipeline
	found := false
	for {
		err := yamlDecoder.Decode(&pipeline)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return nil, err
		}
		if pipeline.IsPipeline() {
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

func (g *GitServiceClient) GetAll() ([]*Pipeline, error) {
	panic("not implemented")
}
