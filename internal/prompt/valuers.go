package prompt

type ValueGetter func() ([]string, error)

var yesNoOptions = []string{"Yes", "No"}

func ListValues(vals []string) ValueGetter {
	return func() ([]string, error) {
		return vals, nil
	}
}

func YesNo() ValueGetter {
	return func() ([]string, error) {
		return yesNoOptions, nil
	}
}
