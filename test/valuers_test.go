package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
)

func TestStaticValidator(t *testing.T) {
	expected := []string{"Phineas", "Ferb"}

	actual, err := prompt.ListValues(expected)()

	assert.Nil(t, err, "expected there to be no errors")
	assert.Contains(t, actual, "Phineas", "expected to find Phineas in the list of values")
	assert.Contains(t, actual, "Ferb", "expected to find Ferb in the list of values")
	assert.NotContains(t, actual, "Perry", "expected to not find Perry in the list of values")
}
