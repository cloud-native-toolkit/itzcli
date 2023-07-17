package prompt

import (
	"fmt"
	"strings"
)

func BaseOptionValidator(p *Prompt, val string) (bool, error) {
	// check to see if the prompt has options, and if it does,
	// then the val must be one of the options.
	if p == nil {
		return false, fmt.Errorf("cannot validate a nil prompt")
	}

	if len(p.AvailableOptions()) == 0 {
		return false, fmt.Errorf("no options defined")
	}

	for _, v := range p.AvailableOptions() {
		if v.value == val {
			return true, nil
		}
	}

	return false, nil
}

func CaseInsensitveTextOptionValidator(p *Prompt, val string) (bool, error) {
	// check to see if the prompt has options, and if it does,
	// then the val must be one of the options.
	if p == nil {
		return false, fmt.Errorf("cannot validate a nil prompt")
	}

	if len(p.AvailableOptions()) == 0 {
		return false, fmt.Errorf("no options defined")
	}
	lval := strings.ToLower(val)
	for _, v := range p.AvailableOptions() {
		if strings.ToLower(v.text) == lval || strings.ToLower(v.Value()) == lval {
			return true, nil
		}
	}

	return false, nil
}
