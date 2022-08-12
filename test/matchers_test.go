package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.ibm.com/Nathan-Good/atkcli/internal/prompt"
)

type MyKey string

func TestAlwaysMatcher(t *testing.T) {

	actual := prompt.Always(context.Background())
	assert.True(t, actual, "the Always matcher should always return true")
}

func TestContextContainsMatcher(t *testing.T) {
	ctx := context.Background()
	actual := prompt.ContextValueIs(ctx, "Skynet", "self-aware")
	assert.False(t, actual, "nothing up our sleves...")

	key := MyKey("Skynet")

	ctx = context.WithValue(ctx, key, "self-aware")
	actual = prompt.ContextValueIs(ctx, key, "self-aware")
	assert.True(t, actual, "expected true when the context does contain the key value")

	actual = prompt.ContextValueIs(ctx, key, "not self-aware")
	assert.False(t, actual, "expected true when the context does contain the key value")
}
