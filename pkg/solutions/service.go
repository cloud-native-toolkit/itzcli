package solutions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"text/tabwriter"
	"text/template"
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

// Solution
type Solution struct {
	Entity *backstage.Entity
}

// Version of the solution gets the version string from the entity's metadata.
// This is a convenience property that pulls the value from the label with the
// value "techzone.ibm.com/version".
func (s *Solution) Version() string {
	val, found := lookupInMap(s.Entity.Metadata.Labels, EntityVersionLabelName)
	if found {
		return val
	}
	return "<unknown>"
}

// PipelineURL gets the pipeline URL from the entity's metadata. This should be
// the full path of the Pipeline YAML file in a GitHub repository.
func (s *Solution) PipelineURL() string {
	val, found := lookupInMap(s.Entity.Metadata.Annotations, PipelineLocationAnnotationName)
	if found {
		return val
	}
	return "<unknown>"
}

type SolutionServiceClient interface {
	Get(id string) (*Solution, error)
	GetAll(f *Filter) ([]Solution, error)
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

func (c WebServiceClient) GetAll(f *Filter) ([]Solution, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("catalog client not properly configured")
	}
	logger.Tracef("Getting entity with filters: %s...", f.BuildFilter())
	e, r, err := c.Client.Catalog.Entities.List(context.Background(), &backstage.ListEntityOptions{
		Filters: f.BuildFilter(),
	})
	if r.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("Error 401 Unauthorized.")
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not get catalog entity with filter: %s", f.BuildFilter()))
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recevied non-success status code: %d", r.StatusCode)
	}
	return toSolutions(e), nil
}

func toSolutions(e []backstage.Entity) []Solution {
	solutions := make([]Solution, len(e))
	for i := range e {
		solutions[i] = Solution{&e[i]}
	}
	return solutions
}

func generateBackStageJWT(c *configuration.ApiConfig) (string, error) {
	techZoneToken := viper.GetString("techzone.api.token")
	if techZoneToken == "" {
		return "", errors.New("No API token set. Please run itz login first, then re-run this command.")
	}
	authEndpoint := fmt.Sprintf("%s/api/rest-login", c.URL)
	backstageToken, err := pkg.ReadHttpGetT(authEndpoint, techZoneToken)
	if err != nil {
		return "", err
	}
	return string(backstageToken), nil
}

func NewWebServiceClient(c *configuration.ApiConfig) (SolutionServiceClient, error) {
	logger.Tracef("Creating backstage client with url: %s...", c.URL)
	bearerToken, err := generateBackStageJWT(c)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: bearerToken,
		TokenType:   "Bearer",
	}))
	client, err := backstage.NewClient(c.URL, DefaultBackstageNamespace, httpClient)
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
	tab := tabwriter.NewWriter(w, 30, 4, 2, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(tab, "NAME\tID\tNAMESPACE\t")
	for _, s := range ss {
		fmt.Fprintln(tab, fmt.Sprintf("%s\t%s\t%s\t", s.Entity.Metadata.Title, s.Entity.Metadata.UID, s.Entity.Metadata.Namespace))
	}
	return tab.Flush()
}

type JsonSolutionWriter struct{}

func (j *JsonSolutionWriter) Write(w io.Writer, s *Solution) error {
	bytes, err := json.Marshal(s.Entity)
	if err == nil {
		w.Write(bytes)
	}
	return err
}

type SolutionInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"UID"`
}

func (j *JsonSolutionWriter) WriteMany(w io.Writer, ss []Solution) error {
	bytes, err := json.Marshal(toSolutionInfo(ss))
	if err == nil {
		w.Write(bytes)
	}
	return err
}

func toSolutionInfo(ss []Solution) []SolutionInfo {
	solutions := make([]SolutionInfo, len(ss))
	for i := range ss {
		solutions[i] = SolutionInfo{
			Name:      ss[i].Entity.Metadata.Title,
			Namespace: ss[i].Entity.Metadata.Namespace,
			UID:       ss[i].Entity.Metadata.UID,
		}
	}
	return solutions
}

func init() {
	writers.Register("text", &TextSolutionWriter{})
	writers.Register("default", &TextSolutionWriter{})
	writers.Register("json", &JsonSolutionWriter{})
}
