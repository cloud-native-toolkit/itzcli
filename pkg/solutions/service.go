package solutions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

var writers SolutionWriters

const DefaultBackstageNamespace string = "default"
const EntityVersionLabelName string = "techzone.ibm.com/version"
const PipelineLocationAnnotationName string = "techzone.ibm.com/tekton-pipeline-location"

func lookupInMap(m map[string]string, key string) (string, bool) {
	for k, v := range m {
		if k == key {
			return v, true
		}
	}
	return "<unknown>", false
}

//go:generate echo "moo"
type Solution struct {
	Entity *backstage.Entity
}

func (s *Solution) Version() string {
	val, found := lookupInMap(s.Entity.Metadata.Labels, EntityVersionLabelName)
	if found {
		return val
	}
	return "<unknown>"
}

func (s *Solution) PipelineURL() string {
	val, found := lookupInMap(s.Entity.Metadata.Annotations, PipelineLocationAnnotationName)
	if found {
		return val
	}
	return "<unknown>"
}

type SolutionServiceClient interface {
	Get(id string) (*Solution, error)
	GetAll(f Filter) ([]Solution, error)
}

type WebServiceClient struct {
	Client  *backstage.Client
	BaseURL string
	Token   string
}

// Get returns the Solution or an error.
func (c WebServiceClient) Get(id string) (*Solution, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("catalog client not properly configured")
	}
	logger.Tracef("Getting entity with id: %s...", id)
	e, r, err := c.Client.Catalog.Entities.Get(context.Background(), id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not get catalog entitry with id %s", id))
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recevied non-success status code: %d", r.StatusCode)
	}
	sol := &Solution{e}
	return sol, nil
}

func (c WebServiceClient) GetAll(f Filter) ([]Solution, error) {
	return nil, nil
}

func NewWebServiceClient(c *configuration.ApiConfig) (SolutionServiceClient, error) {
	logger.Tracef("Creating backstage client with url: %s...", c.URL)
	client, err := backstage.NewClient(c.URL, DefaultBackstageNamespace, nil)
	if err != nil {
		return nil, err
	}
	return &WebServiceClient{
		Client:  client,
		BaseURL: c.URL,
		Token:   c.Token,
	}, nil
}

type SolutionWriter interface {
	Write(w io.Writer, s *Solution) error
	WriteMany(w io.Writer, ss []Solution) error
}

type SolutionWriters struct {
	registered map[string]SolutionWriter
}

func (w *SolutionWriters) Register(name string, writer SolutionWriter) {
	if w.registered == nil {
		w.registered = make(map[string]SolutionWriter)
	}
	w.registered[name] = writer
}

func (w *SolutionWriters) Load(name string) SolutionWriter {
	r := w.registered[name]
	if r == nil {
		return w.registered["default"]
	}
	return r
}

func NewSolutionWriter(format string) SolutionWriter {
	return writers.Load(format)
}

type TextSolutionWriter struct{}

func (t *TextSolutionWriter) Write(w io.Writer, s *Solution) error {
	consoleTemplate := `{{.Entity.Metadata.Title}} 
Version: {{.Version}}
Description: {{.Entity.Metadata.Description}}
Pipeline URL: {{.PipelineURL}}
`
	tmpl, err := template.New("atkrez").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(w, s)
	}
	return nil
}

func (t *TextSolutionWriter) WriteMany(w io.Writer, ss []Solution) error {
	return nil
}

type JsonSolutionWriter struct{}

func (j *JsonSolutionWriter) Write(w io.Writer, s *Solution) error {
	bytes, err := json.Marshal(s.Entity)
	if err == nil {
		w.Write(bytes)
	}
	return err
}

func (j *JsonSolutionWriter) WriteMany(w io.Writer, ss []Solution) error {
	return nil
}

func init() {
	writers.Register("text", &TextSolutionWriter{})
	writers.Register("default", &TextSolutionWriter{})
	writers.Register("json", &JsonSolutionWriter{})
}
