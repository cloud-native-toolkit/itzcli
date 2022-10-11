package test

import (
	"context"
	"fmt"
	"testing"

	"github.ibm.com/skol/atkcli/internal/prompt"

	"github.com/stretchr/testify/assert"
)

func ExistingSubnet(ctx context.Context) bool {
	return ctx.Value("vpc") != nil && ctx.Value("vpc") != "new"
}

/* TestPromptCreateBasic tests basic prompt creation with the builder. */
func TestPromptCreateBasic(t *testing.T) {
	// Create a bunch of prompts with answers to make sure
	// iterators, validators, etc. work correctly.

	builder := prompt.NewPromptBuilder()

	question, err := builder.Path("vpc").
		Text("Would you like to create a new or existing VPC?").
		Build()

	assert.Nil(t, err, "must create without errors")
	assert.NotNil(t, builder, "builder should be not-null")
	assert.Equal(t, fmt.Sprintf("%s", question), "Would you like to create a new or existing VPC?", "String() should print out the prompt text")
	// assert.Equal(t, ctx.Value("vpc"), "vpc", "path should be correct")

}

/* TestSubPrompt provides a test for the PromptSetBuilder, which allows you to build a set of prompts */
func TestSubPrompt(t *testing.T) {
	cloudProviderList := []string{"AWS", "Azure", "GCP"}

	providerQuestion, err := prompt.NewPromptBuilder().
		Path("cloud-provider").
		Text("What cloud provider(s) would you like to use?").
		WithOptions(prompt.ListValues(cloudProviderList)).
		Build()

	assert.Nil(t, err, "must create without errors")

	newOrExistingVpc, err := prompt.NewPromptBuilder().
		Path("vpc").
		Text("First sub question").
		AskWhen(prompt.Always).
		Build()
	assert.Nil(t, err, "must create without errors")

	newOrExistingVpc2, err := prompt.NewPromptBuilder().
		Path("vpc").
		Text("Second sub question").
		AskWhen(prompt.Always).
		Build()
	assert.Nil(t, err, "must create without errors")

	newOrExistingVpc3, err := prompt.NewPromptBuilder().
		Path("vpc").
		Text("Third sub question").
		AskWhen(prompt.Always).
		Build()
	assert.Nil(t, err, "must create without errors")

	providerQuestion.AddSubPrompt(newOrExistingVpc)
	providerQuestion.AddSubPrompt(newOrExistingVpc2)
	providerQuestion.AddSubPrompt(newOrExistingVpc3)

	next := providerQuestion.Itr()

	assert.Equal(t, providerQuestion.String(), next().String(), "the Itr() with a single sub prompt should return the first sub prompt")
	assert.Equal(t, newOrExistingVpc.String(), next().String(), "the Itr() should return the next question")
	assert.Equal(t, newOrExistingVpc2.String(), next().String(), "the Itr() should return the (second) next question")
	assert.Equal(t, newOrExistingVpc3.String(), next().String(), "the Itr() should return the (third) next question")
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

// TestInvalidAnswer tests that if an answer to a prompt is invalid that when next() is called, the prompt is just asked again
func TestLookupAnswer(t *testing.T) {

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
	providerQuestion.Record("AWS")
	assert.Nil(t, next(), "there should be no more because the last answer was correct")

	assert.Equal(t, rootQuestion.GetAnswer("cloud-provider"), "AWS")
}
