package solutions

import (
	"encoding/json"
	"fmt"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
	"io"
	"strings"
	"text/template"
)

type Writer struct {
	jsonFormat bool
}

func (w *Writer) Write(out io.Writer, sols []backstage.Entity) error {
	if w.jsonFormat {
		return JSONWrite(out, sols)
	}
	return TextWriter(out, sols)
}

func NewWriter(jsonFormat bool) *Writer {
	return &Writer{jsonFormat: jsonFormat}
}

func TextWriter(out io.Writer, sols []backstage.Entity) error {
	// TODO: Probably get this from a resource file of some kind
	consoleTemplate := ` - {{.Metadata.Namespace}}/{{.Metadata.Name}} (id: {{.Metadata.UID}})
`
	tmpl, err := template.New("atksol").Parse(consoleTemplate)
	if err == nil {
		for _, sol := range sols {
			err = tmpl.Execute(out, sol)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type solution struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"UID"`
}

func JSONWrite(out io.Writer, sols []backstage.Entity) error {
	var solutions []solution
	for _, sol := range sols {
		s := solution {
			Name:      sol.Metadata.Name,
			Namespace: sol.Metadata.Namespace,
			UID:       sol.Metadata.UID,
		}
		solutions = append(solutions, s)
	}
	bytes, _ := json.Marshal(solutions)
	fmt.Fprint(out, string(bytes))
	return nil
}

type Filter struct {
	Filter []string
}

type FilterOptions func(*Filter)

func OwnerFilter(owner []string) FilterOptions {
	return func(f *Filter) {
		if len(owner) == 0 {
			return
		}
		ownerString := fmt.Sprintf("spec.owner=group:%s", strings.Join(owner, ",spec.owner=group:"))
		f.Filter = append(f.Filter, ownerString)
	}
}

func KindFilter(kind []string) FilterOptions {
	return func(f *Filter) {
		if len(kind) == 0 {
			return
		}
		filterString := fmt.Sprintf("kind=%s", strings.Join(kind, ",kind="))
		f.Filter = append(f.Filter, filterString)
	}
}

func NewFilter(options ...FilterOptions) *Filter {
	filter := &Filter{}

	for _, option := range options {
		option(filter)
	}

	return filter
}

func (f *Filter) BuildFilter() []string {
	var filter []string
	if len(f.Filter) == 0 {
		return filter
	}
	filterString := ""
	for index, filter := range f.Filter {
		if index == 0 {
			filterString = fmt.Sprintf("%s", filter)
			continue
		}
		filterString += fmt.Sprintf(",%s", filter)
	}
	filter = append(filter, filterString)
	return filter
}
