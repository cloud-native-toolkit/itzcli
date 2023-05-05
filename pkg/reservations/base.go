package reservations

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

type Writer struct {
	jsonFormat bool
}

func (w *Writer) Write(out io.Writer, rez TZReservation) error {
	// TODO: Probably get this from a resource file of some kind
	if w.jsonFormat {
		var reservations []TZReservation
		reservations = append(reservations, rez)
		return JSONWrite(out, reservations)
	}
	return TextWriter(out, rez)

}

func TextWriter(out io.Writer, rez TZReservation) error {
	// TODO: Probably get this from a resource file of some kind
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

func JSONWrite(out io.Writer, rez []TZReservation) error {
	jsonData, err := json.Marshal(rez)
	if err != nil {
		return err
	}
	var data TZReservation
	jsonError  := json.Unmarshal(jsonData, &data)
	if jsonError != nil {
		fmt.Println(jsonError)
	}
	fmt.Fprint(out, string(jsonData))
	return nil
}

func (w *Writer) WriteOne(out io.Writer, rez TZReservation) error {
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

func (w *Writer) WriteAll(out io.Writer, rez []TZReservation) error {
	for _, r := range rez {
		err := w.Write(out, r)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (w *Writer) WriteFilter(out io.Writer, rez []TZReservation, filter Filter) (int, error) {
	matches := 0
	var reservations []TZReservation
	for _, r := range rez {
		if filter(r) {
			matches += 1
			// If we need to output as JSON, we need to build an array of filtered reservations
			// so don't call the write function until we have all the reservations 
			if w.jsonFormat {
				reservations = append(reservations, r)
			} else {
				err := w.Write(out, r)
				if err != nil {
					return matches, nil
				}
			}
		}
	}
	if w.jsonFormat {
		err := JSONWrite(out, reservations)
		return matches, err
	}
	return matches, nil
}

func NewWriter(jsonFormat bool) *Writer {
	return &Writer{jsonFormat: jsonFormat}
}
