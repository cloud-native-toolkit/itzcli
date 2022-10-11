package test

import (
	"testing"

	"github.ibm.com/skol/atkcli/internal/prompt"

	"github.com/stretchr/testify/assert"
)

func TestBaseOptionValidator(t *testing.T) {

	question, err := prompt.NewPromptBuilder().
		Path("color").
		Text("What is your favorite color?").
		AddOption("red").
		Build()

	assert.Nil(t, err, "expecting no errors building this question")

	isValid, err := prompt.BaseOptionValidator(question, "red")

	assert.True(t, isValid, "expecting this to be a valid option")
	assert.Nil(t, err, "expecting to not get an error")

	isValid, err = prompt.BaseOptionValidator(question, "blue")
	assert.False(t, isValid, "expecting this to be an invalid option")
	assert.Nil(t, err, "expecting to not get an error even with an invalid option")
}
