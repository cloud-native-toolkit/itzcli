package test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
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

func TestLongAnswer(t *testing.T) {
	rootQuestion, _ := prompt.NewPromptBuilder().
		Path("root").
		Text("Do you want to continue?").
		Build()

	r := strings.NewReader(`eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcElkLWE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3Mi0yMDIxLTAzLTIzVDE0OjQyOjU2LjE5NiIsInZlciI6NH0.eyJpc3MiOiJodHRwczovL2V1LWRlLmFwcGlkLmNsb3VkLmlibS5jb20vb2F1dGgvdjQvYThiZjE3YmMtMTBmNS00NzZiLWE0YzgtYjllYjFlNWQ2MDcyIiwiZXhwIjoxNjY3NDEyMjM2LCJhdWQiOlsiM2E5YzFjZTQtMjk4MC00YzcxLTkxNzEtMmEzY2NjNDNjNzU0Il0sInN1YiI6ImJkYmQwN2MwLWIwOWUtNDEwMi04NWQ1LTc3YWJmZmFhODg4MSIsImFtciI6WyJpYm1pZCJdLCJpYXQiOjE2Njc0MDE0MzYsInRlbmFudCI6ImE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3MiIsInNjb3BlIjoib3BlbmlkIGFwcGlkX2RlZmF1bHQgYXBwaWRfcmVhZHVzZXJhdHRyIGFwcGlkX3JlYWRwcm9maWxlIGFwcGlkX3dyaXRldXNlcmF0dHIgYXBwaWRfYXV0aGVudGljYXRlZCBzdXBlcl9lZGl0IGVkaXQgcmVhZCJ9.cFWkt1_8KtH-4DVfZw9D2gyJv6Bhagw7H3WFXZDIhZvGtH_y6t6VjJCddqnOuboA3ByxBuYreHNA5Kv6d3tD-MHZBLiBsVqkB6qBJ_SXzfRpUeXCvxnm8-O-eVVaazL0VUdyNll_NYuHqxJtrRUTPE-GMhJjmX067mCEv-iiNqEe0v42AQSEwTmBbl2Bfsp7XnYrcT00DohzYv-cbWGMYb2D0KMcCFAdDnile5T5Gz7FRcBzOP7_ZXNSbGdy2Y1JwozD2g2T-OVHW8DLOykkb5h2YdV39Mf1PyLprit5xBAyIZK6Cia9qe1L6WNp8z4QNJrmEVeL6_5bP6K04PyazA eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcElkLWE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3Mi0yMDIxLTAzLTIzVDE0OjQyOjU2LjE5NiIsInZlciI6NH0.eyJpc3MiOiJodHRwczovL2V1LWRlLmFwcGlkLmNsb3VkLmlibS5jb20vb2F1dGgvdjQvYThiZjE3YmMtMTBmNS00NzZiLWE0YzgtYjllYjFlNWQ2MDcyIiwiYXVkIjpbIjNhOWMxY2U0LTI5ODAtNGM3MS05MTcxLTJhM2NjYzQzYzc1NCJdLCJleHAiOjE2Njc0MTIyMzYsInRlbmFudCI6ImE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3MiIsImlhdCI6MTY2NzQwMTQzNiwiZW1haWwiOiJuYXRoYW4uZ29vZEBpYm0uY29tIiwibmFtZSI6Ik5hdGhhbiBHb29kIiwic3ViIjoiYmRiZDA3YzAtYjA5ZS00MTAyLTg1ZDUtNzdhYmZmYWE4ODgxIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiTmF0aGFuLkdvb2RAaWJtLmNvbSIsImdpdmVuX25hbWUiOiJOYXRoYW4iLCJmYW1pbHlfbmFtZSI6Ikdvb2QiLCJpZGVudGl0aWVzIjpbeyJwcm92aWRlciI6ImlibWlkIiwiaWQiOiJuYXRoYW4uZ29vZEBpYm0uY29tIn1dLCJhbXIiOlsiaWJtaWQiXX0.Zz2gmW0NXMPCxZ9CmLKFV0gmxdeCKt8SPKdXwg03x2c98v1mm0eWTfF2hw6APPuZZ8BVdJI8yiwd7hlAxWpIk6dwtBC-c1UkyG23kX4iNEyVnW8XXZqpp-BZ7j88frGWdLb_R-5tsnh0GACofh6Fbp8koSphmN5OPunPRriaZAGZnESJvdiuQ63dmTXC3nU_LbaCROUyzxok2o4ohEoHRSA7-Pf9YPQOsBz5IyFJHdOZreCoE1CRqzKlsd6KLeKXQ2XbsidhTzn_oBoM-hK55dQEvr4rtyUTn1WpRRaFu7kqTKd0Xt85UAAA7Hmkpe_Wv0FqV65KN_xFJl3VBTuF7w`)
	w := bytes.NewBufferString("Do you want to continue?")

	prompt.Ask(rootQuestion, w, r)

	assert.Equal(t, rootQuestion.GetAnswer("root"), `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcElkLWE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3Mi0yMDIxLTAzLTIzVDE0OjQyOjU2LjE5NiIsInZlciI6NH0.eyJpc3MiOiJodHRwczovL2V1LWRlLmFwcGlkLmNsb3VkLmlibS5jb20vb2F1dGgvdjQvYThiZjE3YmMtMTBmNS00NzZiLWE0YzgtYjllYjFlNWQ2MDcyIiwiZXhwIjoxNjY3NDEyMjM2LCJhdWQiOlsiM2E5YzFjZTQtMjk4MC00YzcxLTkxNzEtMmEzY2NjNDNjNzU0Il0sInN1YiI6ImJkYmQwN2MwLWIwOWUtNDEwMi04NWQ1LTc3YWJmZmFhODg4MSIsImFtciI6WyJpYm1pZCJdLCJpYXQiOjE2Njc0MDE0MzYsInRlbmFudCI6ImE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3MiIsInNjb3BlIjoib3BlbmlkIGFwcGlkX2RlZmF1bHQgYXBwaWRfcmVhZHVzZXJhdHRyIGFwcGlkX3JlYWRwcm9maWxlIGFwcGlkX3dyaXRldXNlcmF0dHIgYXBwaWRfYXV0aGVudGljYXRlZCBzdXBlcl9lZGl0IGVkaXQgcmVhZCJ9.cFWkt1_8KtH-4DVfZw9D2gyJv6Bhagw7H3WFXZDIhZvGtH_y6t6VjJCddqnOuboA3ByxBuYreHNA5Kv6d3tD-MHZBLiBsVqkB6qBJ_SXzfRpUeXCvxnm8-O-eVVaazL0VUdyNll_NYuHqxJtrRUTPE-GMhJjmX067mCEv-iiNqEe0v42AQSEwTmBbl2Bfsp7XnYrcT00DohzYv-cbWGMYb2D0KMcCFAdDnile5T5Gz7FRcBzOP7_ZXNSbGdy2Y1JwozD2g2T-OVHW8DLOykkb5h2YdV39Mf1PyLprit5xBAyIZK6Cia9qe1L6WNp8z4QNJrmEVeL6_5bP6K04PyazA eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcElkLWE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3Mi0yMDIxLTAzLTIzVDE0OjQyOjU2LjE5NiIsInZlciI6NH0.eyJpc3MiOiJodHRwczovL2V1LWRlLmFwcGlkLmNsb3VkLmlibS5jb20vb2F1dGgvdjQvYThiZjE3YmMtMTBmNS00NzZiLWE0YzgtYjllYjFlNWQ2MDcyIiwiYXVkIjpbIjNhOWMxY2U0LTI5ODAtNGM3MS05MTcxLTJhM2NjYzQzYzc1NCJdLCJleHAiOjE2Njc0MTIyMzYsInRlbmFudCI6ImE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3MiIsImlhdCI6MTY2NzQwMTQzNiwiZW1haWwiOiJuYXRoYW4uZ29vZEBpYm0uY29tIiwibmFtZSI6Ik5hdGhhbiBHb29kIiwic3ViIjoiYmRiZDA3YzAtYjA5ZS00MTAyLTg1ZDUtNzdhYmZmYWE4ODgxIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiTmF0aGFuLkdvb2RAaWJtLmNvbSIsImdpdmVuX25hbWUiOiJOYXRoYW4iLCJmYW1pbHlfbmFtZSI6Ikdvb2QiLCJpZGVudGl0aWVzIjpbeyJwcm92aWRlciI6ImlibWlkIiwiaWQiOiJuYXRoYW4uZ29vZEBpYm0uY29tIn1dLCJhbXIiOlsiaWJtaWQiXX0.Zz2gmW0NXMPCxZ9CmLKFV0gmxdeCKt8SPKdXwg03x2c98v1mm0eWTfF2hw6APPuZZ8BVdJI8yiwd7hlAxWpIk6dwtBC-c1UkyG23kX4iNEyVnW8XXZqpp-BZ7j88frGWdLb_R-5tsnh0GACofh6Fbp8koSphmN5OPunPRriaZAGZnESJvdiuQ63dmTXC3nU_LbaCROUyzxok2o4ohEoHRSA7-Pf9YPQOsBz5IyFJHdOZreCoE1CRqzKlsd6KLeKXQ2XbsidhTzn_oBoM-hK55dQEvr4rtyUTn1WpRRaFu7kqTKd0Xt85UAAA7Hmkpe_Wv0FqV65KN_xFJl3VBTuF7w`)
}
