package reservations

import (
	"encoding/json"
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

type TZReservation struct {
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

type Filter func(TZReservation) bool

func FilterByStatus(status string) Filter {
	return func(r TZReservation) bool {
		return r.Status == status
	}
}

func FilterByStatusSlice(status []string) Filter {
	return func(r TZReservation) bool {
		return pkg.StringSliceContains(status, r.Status)
	}
}


type OutputWriter interface {
	io.Writer
}

type Reader interface {
	Read(io.Reader) (TZReservation, error)
	ReadAll(io.Reader) ([]TZReservation, error)
}

type Writer interface {
	Write(io.Writer, TZReservation) error
	WriteAll(io.Writer, []TZReservation) error
	WriteFilter(io.Writer, []TZReservation, Filter) error
}

type JsonReader struct{}

func (j *JsonReader) Read(reader io.Reader) (TZReservation, error) {
	var res TZReservation
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func (j *JsonReader) ReadAll(reader io.Reader) ([]TZReservation, error) {
	var res []TZReservation
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func NewJsonReader() *JsonReader {
	return &JsonReader{}
}

type TextWriter struct{}

func (w *TextWriter) Write(out io.Writer, rez TZReservation) error {
	// TODO: Probably get this from a resource file of some kind
	consoleTemplate := ` - {{.Name}} - {{.Status}}
   Reservation Id: {{.ReservationId}}

`
	tmpl, err := template.New("atkrez").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(out, rez)
	}
	return nil
}

func (w *TextWriter) WriteOne(out io.Writer, rez TZReservation) error {
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

func (w *TextWriter) WriteAll(out io.Writer, rez []TZReservation) error {
	for _, r := range rez {
		err := w.Write(out, r)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (w *TextWriter) WriteFilter(out io.Writer, rez []TZReservation, filter Filter) (int, error) {
	matches := 0
	for _, r := range rez {
		if filter(r) {
			matches += 1
			err := w.Write(out, r)
			if err != nil {
				return matches, nil
			}
		}
	}
	return matches, nil
}

func NewTextWriter() *TextWriter {
	return &TextWriter{}
}
