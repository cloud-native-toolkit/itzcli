package solutions

import (
	"encoding/json"
	"io"
	"text/template"
	"github.com/cloud-native-toolkit/itzcli/pkg"
)

type Solution struct {
	Kind string `json:"kind"`
	Metadata `json:"metadata"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}

type Filter func(Solution) bool

func FilterByStatusSlice(status []string) Filter {
	return func(s Solution) bool {
		return pkg.StringSliceContains(status, s.Kind)
	}
}
type Reader interface {
	Read(io.Reader) (Solution, error)
	ReadAll(io.Reader) ([]Solution, error)
}

type Writer interface {
	Write(io.Writer, Solution) error
	WriteAll(io.Writer, []Solution) error
	WriteFilter(io.Writer, []Solution, Filter) error
}

type JsonReader struct{}

func (j *JsonReader) Read(reader io.Reader) (Solution, error) {
	var res Solution
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func (j *JsonReader) ReadAll(reader io.Reader) ([]Solution, error) {
	var res []Solution
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func NewJsonReader() *JsonReader {
	return &JsonReader{}
}

type TextWriter struct{}

func (w *TextWriter) Write(out io.Writer, sol Solution) error {
	// TODO: Probably get this from a resource file of some kind
	consoleTemplate := ` - {{.Metadata.Namespace}}/{{.Metadata.Name}} (id: {{.Metadata.UID}})
`
	tmpl, err := template.New("atksol").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(out, sol)
	}
	return nil
}

func (w *TextWriter) WriteAll(out io.Writer, sol []Solution) error {
	for _, r := range sol {
		err := w.Write(out, r)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (w *TextWriter) WriteFilter(out io.Writer, sols []Solution, filter Filter) error {
	for _, s := range sols {
		if filter(s) {
			err := w.Write(out, s)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

func NewTextWriter() *TextWriter {
	return &TextWriter{}
}
