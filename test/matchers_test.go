package test

import (
	"github.ibm.com/skol/atkcli/internal/prompt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyKey string

func TestAlwaysMatcher(t *testing.T) {
	p, _ := prompt.NewPromptBuilder().Build()
	actual := prompt.Always(p)
	assert.True(t, actual, "the Always matcher should always return true")
}

func TestContextContainsMatcher(t *testing.T) {
	p, err := prompt.NewPromptBuilder().Path("moo").Build()
	assert.Nil(t, err, "should be able to create the prompt")
	actual := prompt.AnsweredValueIs(p, "Skynet", "self-aware")
	assert.False(t, actual, "nothing up our sleves...")

	key := MyKey("Skynet")
	p, _ = prompt.NewPromptBuilder().Path(string(key)).Build()

	p.Record("self-aware")
	actual = prompt.AnsweredValueIs(p, string(key), "self-aware")
	assert.True(t, actual, "expected true when the context does contain the key value")

	p2, _ := prompt.NewPromptBuilder().Path("oink").Build()
	actual = prompt.AnsweredValueIs(p2, string(key), "not self-aware")
	assert.False(t, actual, "expected false when the context does contain the key value")
}
