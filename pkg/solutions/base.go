package solutions

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"go.einride.tech/backstage/catalog"
)



type TextWriter struct{}

func (w *TextWriter) Write(out io.Writer, sol *catalog.Entity) error {
	// TODO: Probably get this from a resource file of some kind
	consoleTemplate := ` - {{.Metadata.Namespace}}/{{.Metadata.Name}} (id: {{.Metadata.UID}})
`
	tmpl, err := template.New("atksol").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(out, sol)
	}
	return nil
}

func NewTextWriter() *TextWriter {
	return &TextWriter{}
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
		ownerString := fmt.Sprintf("spec.owner=group:%s", strings.Join(owner,",spec.owner=group:"))
		f.Filter = append(f.Filter, ownerString)
	}
}

func KindFilter(kind []string) FilterOptions {
	return func(f *Filter) {
		if len(kind) == 0 {
			return
		}
		filterString := fmt.Sprintf("kind=%s", strings.Join(kind,",kind="))
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
	for index, filter := range f.Filter {
		if index != 0 {
			f.Filter[index] = fmt.Sprintf("&%s", filter)
		}
	}
	return f.Filter
}