package solutions

import (
	"encoding/json"
	"io"
	"text/template"
)

type Solution struct {
	Id        string
	Name      string
	Platform  string
	ShortDesc string `json:"short_desc"`
	LongDesc  string `json:"long_desc"`
	Public    bool
}

type Filter func(Solution) bool

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
	consoleTemplate := ` - {{.Name}} (id: {{.Id}})
`
	tmpl, err := template.New("atksol").Parse(consoleTemplate)
	if err == nil {
		tmpl.Execute(out, sol)
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
