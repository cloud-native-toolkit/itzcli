package test

import (
	"testing"

	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"github.com/stretchr/testify/assert"
)

func TestStaticValidator(t *testing.T) {
	expected := []string{"Phineas", "Ferb"}

	actual, err := prompt.ListBasicValues(expected)()

	assert.Nil(t, err, "expected there to be no errors")
	assert.Contains(t, actual, "phineas", "expected to find Phineas in the list of values")
	assert.Contains(t, actual, "ferb", "expected to find Ferb in the list of values")
	assert.NotContains(t, actual, "perry", "expected to not find Perry in the list of values")
}
