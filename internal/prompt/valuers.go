package prompt

type ValueGetter func() ([]string, error)

func ListValues(vals []string) ValueGetter {
	return func() ([]string, error) {
		return vals, nil
	}
}
