package techzone

import (
	"bytes"
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io"
	"reflect"
	"text/template"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
)

var writers RegisteredModelWriters

const DefaultOutputFormat = "text"

type Environment struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ReservationServiceClient interface {
	Get(id string) (*Reservation, error)
	GetAll(f Filter) ([]Reservation, error)
}

type ReservationWebServiceClient struct {
	BaseURL string
	Token   string
}

// Get
func (c *ReservationWebServiceClient) Get(id string) (*Reservation, error) {
	fullUrl := fmt.Sprintf("%s/reservation/ibmcloud-2/%s", c.BaseURL, id)

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of reservations...",
		c.BaseURL, c.Token)

	data, err := pkg.ReadHttpGetTWithFunc(fullUrl, c.Token, func(code int) error {
		logger.Debugf("Handling HTTP return code %d...", code)
		if code == 401 {
			return fmt.Errorf("not authorized")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	jsoner := NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.Read(dataR)
	return &rez, err
}

// GetAll
func (c *ReservationWebServiceClient) GetAll(f Filter) ([]Reservation, error) {
	fullUrl := fmt.Sprintf("%s/my/reservations/all", c.BaseURL)

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of reservations...",
		c.BaseURL, c.Token)

	data, err := pkg.ReadHttpGetTWithFunc(fullUrl, c.Token, func(code int) error {
		logger.Debugf("Handling HTTP return code %d...", code)
		if code == 401 {
			return fmt.Errorf("not authorized")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	jsoner := NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.ReadAll(dataR)
	return rez, err
}

func NewReservationWebServiceClient(c *configuration.ApiConfig) (ReservationServiceClient, error) {
	return &ReservationWebServiceClient{
		BaseURL: c.URL,
		Token:   c.Token,
	}, nil
}

// EnvironmentServiceClient the client API for EnvironmentService service.
type EnvironmentServiceClient interface {
	Get(id string) (*Environment, error)
	GetAll(f Filter) ([]Environment, error)
}

type EnvironmentWebServiceClient struct {
	BaseURL string
	Token   string
}

// Get
func (c *EnvironmentWebServiceClient) Get(id string) (*Environment, error) {
	return nil, nil
}

// GetAll
func (c *EnvironmentWebServiceClient) GetAll(f Filter) ([]Environment, error) {
	return nil, nil
}

func NewEnvironmentWebServiceClient(c *configuration.ApiConfig) (EnvironmentServiceClient, error) {
	return &EnvironmentWebServiceClient{
		BaseURL: c.URL,
		Token:   c.Token,
	}, nil
}

type ModelWriter interface {
	WriteOne(w io.Writer, val interface{}) error
	WriteMany(w io.Writer, val interface{}) error
}

type WriterKey struct {
	modelType    string
	outputFormat string
}

func defaultKey(key WriterKey) WriterKey {
	return WriterKey{
		modelType:    key.modelType,
		outputFormat: DefaultOutputFormat,
	}
}

type RegisteredModelWriters struct {
	registered map[WriterKey]ModelWriter
}

type TextReservationWriter struct{}

func (t *TextReservationWriter) WriteOne(w io.Writer, val interface{}) error {
	// TODO: Probably get this from a resource file of some kind
	consoleTemplate := ` - {{.Name}} - {{.Status}}
   Reservation Id: {{.ReservationId}}
   Description: {{.Description}}
   Collection Id: {{.CollectionId}}
   Extend Count: {{.ExtendCount}}
   Service Links:
    --------------------------------
    {{- range .ServiceLinks}}
		{{- if .Sensitive}}
			{{- printf "\n    %s: ****Private****\n    --------------------------------" .Label}}
		{{- else}} 
			{{- printf "\n    %s: %s\n    --------------------------------" .Label .Url}}
		{{- end}}
	{{- end}}
`

	tmpl, err := template.New("atkrez").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(w, val)
	}
	return nil
}

func (t *TextReservationWriter) WriteMany(w io.Writer, val interface{}) error {
	// TODO: Probably get this from a resource file of some kind
	consoleTemplate := `{{- range .}} - {{.Name}} - {{.Status}}
   Reservation Id: {{.ReservationId}}

{{ end}}`
	tmpl, err := template.New("atkrez").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(w, val)
	}
	return nil
}

type JsonReservationWriter struct{}

func (j *JsonReservationWriter) WriteOne(w io.Writer, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err == nil {
		w.Write(bytes)
	}
	return err
}

func (j *JsonReservationWriter) WriteMany(w io.Writer, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err == nil {
		w.Write(bytes)
	}
	return err
}

func (w *RegisteredModelWriters) Register(forType string, format string, writer ModelWriter) {
	if w.registered == nil {
		w.registered = make(map[WriterKey]ModelWriter)
	}
	key := WriterKey{modelType: forType, outputFormat: format}
	w.registered[key] = writer
}

func (w *RegisteredModelWriters) Load(forType string, format string) ModelWriter {
	key := WriterKey{modelType: forType, outputFormat: format}
	r := w.registered[key]
	if r == nil {
		d := defaultKey(key)
		return w.registered[d]
	}
	return r
}

func NewModelWriter(forType string, format string) ModelWriter {
	return writers.Load(forType, format)
}

func init() {
	reservationType := reflect.TypeOf(Reservation{})
	logger.Tracef("Registering writers for type: %s", reservationType)
	writers.Register(reservationType.Name(), "text", &TextReservationWriter{})
	writers.Register(reservationType.Name(), DefaultOutputFormat, &TextReservationWriter{})
	writers.Register(reservationType.Name(), "json", &JsonReservationWriter{})
}
