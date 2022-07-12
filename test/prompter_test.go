package test

import (
	"context"
	"fmt"
	"testing"

	"github.ibm.com/Nathan-Good/opcnow-go/internal/prompt"

	"github.com/stretchr/testify/assert"
)

func ExistingSubnet(ctx context.Context) bool {
	return ctx.Value("vpc") != nil && ctx.Value("vpc") != "new"
}

/* TestPromptCreateBasic tests basic prompt creation with the builder. */
func TestPromptCreateBasic(t *testing.T) {
	// Create a bunch of prompts with answers to make sure
	// iterators, validators, etc. work correctly.
	ctx := context.Background()

	builder := prompt.NewPromptBuilder()

	question, err := builder.Path("vpc").
		Context(ctx).
		Text("Would you like to create a new or existing VPC?").
		Build()

	assert.Nil(t, err, "must create without errors")
	assert.NotNil(t, builder, "builder should be not-null")
	assert.Equal(t, fmt.Sprintf("%s", question), "Would you like to create a new or existing VPC?", "String() should print out the prompt text")
	// assert.Equal(t, ctx.Value("vpc"), "vpc", "path should be correct")

}

/* TestSubPrompt provides a test for the PromptSetBuilder, which allows you to build a set of prompts */
func TestSubPrompt(t *testing.T) {

	ctx := context.Background()

	cloudProviderList := []string{"AWS", "Azure", "GCP"}

	providerQuestion, err := prompt.NewPromptBuilder().
		Context(ctx).
		Path("cloud-provider").
		Text("What cloud provider(s) would you like to use?").
		WithOptions(prompt.ListValues(cloudProviderList)).
		Build()

	assert.Nil(t, err, "must create without errors")

	newOrExistingVpc, err := prompt.NewPromptBuilder().
		Path("vpc").
		AskWhen(prompt.Always).
		Build()
	assert.Nil(t, err, "must create without errors")

	providerQuestion.AddSubPrompt(newOrExistingVpc)

	next := providerQuestion.Itr()

	assert.Equal(t, providerQuestion.String(), next().String(), "the Itr() with a single sub prompt should return the first sub prompt")
	assert.Equal(t, newOrExistingVpc.String(), next().String(), "the Itr() shoud return the next question")
	assert.Nil(t, next(), "Should be able to safely call next() when none are left and get nil")

}

// TestInvalidAnswer tests that if an answer to a prompt is invalid that when next() is called, the prompt is just asked again
func TestInvalidAnswer(t *testing.T) {

	cloudProviderList := []string{"AWS", "Azure", "GCP"}

	rootQuestion, _ := prompt.NewPromptBuilder().
		Path("root").
		Text("Do you want to continue?").
		Build()

	providerQuestion, _ := prompt.NewPromptBuilder().
		Path("cloud-provider").
		Text("What cloud provider(s) would you like to use?").
		WithOptions(prompt.ListValues(cloudProviderList)).
		WithValidator(prompt.BaseOptionValidator).
		Build()

	rootQuestion.AddSubPrompt(providerQuestion)

	next := rootQuestion.Itr()

	assert.Equal(t, rootQuestion.String(), next().String(), "the first root question")
	assert.Equal(t, providerQuestion.String(), next().String(), "the first sub menu item")
	providerQuestion.Record("Moo")
	assert.Equal(t, providerQuestion.String(), next().String(), "incorrect answer, ask again")
	providerQuestion.Record("AWS")
	assert.Nil(t, next(), "there should be no more because the last answer was correct")
}
