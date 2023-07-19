package solutions

import (
	"fmt"
	"strings"
)

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

func ComponentNameFilter(name []string) FilterOptions {
	return func(f *Filter) {
		if len(name) == 0 {
			return
		}
		filterString := fmt.Sprintf("metadata.title=%s", strings.Join(name, ",metadata.title="))
		f.Filter = append(f.Filter, filterString)
	}
}

func TypeFilter(t []string) FilterOptions {
	return func(f *Filter) {
		if len(t) == 0 {
			return
		}
		filterString := fmt.Sprintf("spec.type=%s", strings.Join(t, ",spec.type="))
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
