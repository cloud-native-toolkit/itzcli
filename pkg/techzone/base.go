package techzone

import (
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/cloud-native-toolkit/itzcli/pkg"
)

type ServiceLink struct {
	LinkType  string `json:"type"`
	Label     string
	Sensitive bool
	Url       string
}

type Reservation struct {
	Name           string
	ServiceLinks   []ServiceLink
	OpportunityId  []string
	ReservationId  string `json:"id"`
	CreatedAt      int
	Status         string
	ProvisionDate  string
	ProvisionUntil string
	CollectionId   string
	ExtendCount    int
	Description    string
}

type Filter func(Reservation) bool

func NoFilter() Filter {
	return func(r Reservation) bool {
		return true
	}
}

func FilterByStatus(status string) Filter {
	return func(r Reservation) bool {
		return r.Status == status
	}
}

func FilterByStatusSlice(status []string) Filter {
	return func(r Reservation) bool {
		return pkg.StringSliceContains(status, r.Status)
	}
}

type OutputWriter interface {
	io.Writer
}

type Reader interface {
	Read(io.Reader) (Reservation, error)
	ReadAll(io.Reader) ([]Reservation, error)
}

type JsonReader struct{}

func (j *JsonReader) Read(reader io.Reader) (Reservation, error) {
	var res Reservation
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func (j *JsonReader) ReadAll(reader io.Reader) ([]Reservation, error) {
	var res []Reservation
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func NewJsonReader() *JsonReader {
	return &JsonReader{}
}

type ReservationTextWriter struct{}

func (w *ReservationTextWriter) WriteOne(out io.Writer, rez *Reservation) error {
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
		return tmpl.Execute(out, rez)
	}
	return nil
}

func (w *ReservationTextWriter) WriteMany(out io.Writer, rez []Reservation) error {
	// TODO: Probably get this from a resource file of some kind
	consoleTemplate := `{{- range .}} - {{.Name}} - {{.Status}}
   Reservation Id: {{.ReservationId}}

{{ end}}`
	tmpl, err := template.New("atkrez").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(out, rez)
	}
	return nil
}

type ReservationJsonWriter struct{}

func (w ReservationJsonWriter) WriteOne(out io.Writer, rez *Reservation) error {
	jsonData, err := json.Marshal(rez)
	if err != nil {
		return err
	}
	b, err := out.Write(jsonData)
	if b == 0 {
		return fmt.Errorf("unexpected writing zero bytes")
	}
	return err
}

func (w ReservationJsonWriter) WriteMany(out io.Writer, rez []Reservation) error {
	jsonData, err := json.Marshal(rez)
	if err != nil {
		return err
	}
	b, err := out.Write(jsonData)
	if b == 0 {
		return fmt.Errorf("unexpected writing zero bytes")
	}
	return err
}

func WriteReservation(w ReservationWriter, out io.Writer, rez *Reservation) error {
	return w.WriteOne(out, rez)
}

// WriteFilteredReservations writes one or more reservations, if they pass the filter. To
// print all of them (basically without filtering, use NoFilter
func WriteFilteredReservations(w ReservationWriter, out io.Writer, rez []Reservation, filter Filter) (int, error) {
	matches := 0
	var filtered []Reservation
	for _, r := range rez {
		if filter(r) {
			matches += 1
			filtered = append(filtered, r)
		}
	}
	err := w.WriteMany(out, filtered)
	return matches, err
}

func NewWriter(format string) ReservationWriter {
	if format == "json" {
		return &ReservationJsonWriter{}
	}
	return &ReservationTextWriter{}
}
