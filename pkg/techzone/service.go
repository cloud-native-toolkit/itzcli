package techzone

import (
	"io"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
)

var reservationWriters ReservationWriters

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
	return nil, nil
}

// GetAll
func (c *ReservationWebServiceClient) GetAll(f Filter) ([]Reservation, error) {
	return nil, nil
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

type EnvironmentWriter interface {
	WriteOne(w io.Writer, s *Environment) error
	WriteMany(w io.Writer, ss []Environment) error
}

type ReservationWriter interface {
	WriteOne(w io.Writer, s *Reservation) error
	WriteMany(w io.Writer, ss []Reservation) error
}

type ReservationWriters struct {
	registered map[string]ReservationWriter
}

type TextReservationWriter struct{}

func (t *TextReservationWriter) WriteOne(w io.Writer, s *Reservation) error {
	return nil
}

func (t *TextReservationWriter) WriteMany(w io.Writer, ss []Reservation) error {
	return nil
}

type JsonReservationWriter struct{}

func (w *ReservationWriters) Register(name string, writer ReservationWriter) {
	if w.registered == nil {
		w.registered = make(map[string]ReservationWriter)
	}
	w.registered[name] = writer
}

func (w *ReservationWriters) Load(name string) ReservationWriter {
	r := w.registered[name]
	if r == nil {
		return w.registered["default"]
	}
	return r
}

func NewReservationWriter(format string) ReservationWriter {
	return reservationWriters.Load(format)
}

func init() {
	reservationWriters.Register("text", &TextReservationWriter{})
	reservationWriters.Register("default", &TextReservationWriter{})
	reservationWriters.Register("json", &JsonReservationWriter{})
}
