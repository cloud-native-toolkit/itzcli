package prompt

import "strings"

type ValueGetter func() (map[string]string, error)

var yesNoOptions = []string{"Yes", "No"}

func ListBasicValues(vals []string) ValueGetter {
	return func() (map[string]string, error) {
		result := make(map[string]string, 0)
		for _, v := range vals {
			result[strings.ToLower(v)] = v
		}
		return result, nil
	}
}

func YesNo() ValueGetter {
	return func() (map[string]string, error) {
		return ListBasicValues(yesNoOptions)()
	}
}
