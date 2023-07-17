package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"reflect"
	"strings"
	"testing"

	"github.com/cloud-native-toolkit/itzcli/internal/prompt"

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
		WithOptions(prompt.ListBasicValues(cloudProviderList)).
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
		WithOptions(prompt.ListBasicValues(cloudProviderList)).
		WithValidator(prompt.BaseOptionValidator).
		Build()

	rootQuestion.AddSubPrompt(providerQuestion)

	next := rootQuestion.Itr()

	assert.Equal(t, rootQuestion.String(), next().String(), "the first root question")
	assert.Equal(t, providerQuestion.String(), next().String(), "the first sub menu item")
	providerQuestion.Record("Moo")
	assert.Equal(t, providerQuestion.String(), next().String(), "incorrect answer, ask again")
	providerQuestion.Record("aws")
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
		WithOptions(prompt.ListBasicValues(cloudProviderList)).
		WithValidator(prompt.BaseOptionValidator).
		Build()

	rootQuestion.AddSubPrompt(providerQuestion)

	next := rootQuestion.Itr()

	assert.Equal(t, rootQuestion.String(), next().String(), "the first root question")
	assert.Equal(t, providerQuestion.String(), next().String(), "the first sub menu item")
	providerQuestion.Record("aws")
	assert.Nil(t, next(), "there should be no more because the last answer was correct")

	assert.Equal(t, rootQuestion.GetAnswer("cloud-provider"), "aws")
}

func TestLongAnswer(t *testing.T) {
	rootQuestion, _ := prompt.NewPromptBuilder().
		Path("root").
		Text("Do you want to continue?").
		Build()

	r := strings.NewReader(`eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcElkLWE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3Mi0yMDIxLTAzLTIzVDE0OjQyOjU2LjE5NiIsInZlciI6NH0.eyJpc3MiOiJodHRwczovL2V1LWRlLmFwcGlkLmNsb3VkLmlibS5jb20vb2F1dGgvdjQvYThiZjE3YmMtMTBmNS00NzZiLWE0YzgtYjllYjFlNWQ2MDcyIiwiZXhwIjoxNjY3NDEyMjM2LCJhdWQiOlsiM2E5YzFjZTQtMjk4MC00YzcxLTkxNzEtMmEzY2NjNDNjNzU0Il0sInN1YiI6ImJkYmQwN2MwLWIwOWUtNDEwMi04NWQ1LTc3YWJmZmFhODg4MSIsImFtciI6WyJpYm1pZCJdLCJpYXQiOjE2Njc0MDE0MzYsInRlbmFudCI6ImE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3MiIsInNjb3BlIjoib3BlbmlkIGFwcGlkX2RlZmF1bHQgYXBwaWRfcmVhZHVzZXJhdHRyIGFwcGlkX3JlYWRwcm9maWxlIGFwcGlkX3dyaXRldXNlcmF0dHIgYXBwaWRfYXV0aGVudGljYXRlZCBzdXBlcl9lZGl0IGVkaXQgcmVhZCJ9.cFWkt1_8KtH-4DVfZw9D2gyJv6Bhagw7H3WFXZDIhZvGtH_y6t6VjJCddqnOuboA3ByxBuYreHNA5Kv6d3tD-MHZBLiBsVqkB6qBJ_SXzfRpUeXCvxnm8-O-eVVaazL0VUdyNll_NYuHqxJtrRUTPE-GMhJjmX067mCEv-iiNqEe0v42AQSEwTmBbl2Bfsp7XnYrcT00DohzYv-cbWGMYb2D0KMcCFAdDnile5T5Gz7FRcBzOP7_ZXNSbGdy2Y1JwozD2g2T-OVHW8DLOykkb5h2YdV39Mf1PyLprit5xBAyIZK6Cia9qe1L6WNp8z4QNJrmEVeL6_5bP6K04PyazA eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcElkLWE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3Mi0yMDIxLTAzLTIzVDE0OjQyOjU2LjE5NiIsInZlciI6NH0.eyJpc3MiOiJodHRwczovL2V1LWRlLmFwcGlkLmNsb3VkLmlibS5jb20vb2F1dGgvdjQvYThiZjE3YmMtMTBmNS00NzZiLWE0YzgtYjllYjFlNWQ2MDcyIiwiYXVkIjpbIjNhOWMxY2U0LTI5ODAtNGM3MS05MTcxLTJhM2NjYzQzYzc1NCJdLCJleHAiOjE2Njc0MTIyMzYsInRlbmFudCI6ImE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3MiIsImlhdCI6MTY2NzQwMTQzNiwiZW1haWwiOiJuYXRoYW4uZ29vZEBpYm0uY29tIiwibmFtZSI6Ik5hdGhhbiBHb29kIiwic3ViIjoiYmRiZDA3YzAtYjA5ZS00MTAyLTg1ZDUtNzdhYmZmYWE4ODgxIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiTmF0aGFuLkdvb2RAaWJtLmNvbSIsImdpdmVuX25hbWUiOiJOYXRoYW4iLCJmYW1pbHlfbmFtZSI6Ikdvb2QiLCJpZGVudGl0aWVzIjpbeyJwcm92aWRlciI6ImlibWlkIiwiaWQiOiJuYXRoYW4uZ29vZEBpYm0uY29tIn1dLCJhbXIiOlsiaWJtaWQiXX0.Zz2gmW0NXMPCxZ9CmLKFV0gmxdeCKt8SPKdXwg03x2c98v1mm0eWTfF2hw6APPuZZ8BVdJI8yiwd7hlAxWpIk6dwtBC-c1UkyG23kX4iNEyVnW8XXZqpp-BZ7j88frGWdLb_R-5tsnh0GACofh6Fbp8koSphmN5OPunPRriaZAGZnESJvdiuQ63dmTXC3nU_LbaCROUyzxok2o4ohEoHRSA7-Pf9YPQOsBz5IyFJHdOZreCoE1CRqzKlsd6KLeKXQ2XbsidhTzn_oBoM-hK55dQEvr4rtyUTn1WpRRaFu7kqTKd0Xt85UAAA7Hmkpe_Wv0FqV65KN_xFJl3VBTuF7w`)
	w := bytes.NewBufferString("Do you want to continue?")

	err := prompt.Ask(rootQuestion, w, r)
	assert.NoError(t, err)
	assert.Equal(t, rootQuestion.GetAnswer("root"), `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcElkLWE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3Mi0yMDIxLTAzLTIzVDE0OjQyOjU2LjE5NiIsInZlciI6NH0.eyJpc3MiOiJodHRwczovL2V1LWRlLmFwcGlkLmNsb3VkLmlibS5jb20vb2F1dGgvdjQvYThiZjE3YmMtMTBmNS00NzZiLWE0YzgtYjllYjFlNWQ2MDcyIiwiZXhwIjoxNjY3NDEyMjM2LCJhdWQiOlsiM2E5YzFjZTQtMjk4MC00YzcxLTkxNzEtMmEzY2NjNDNjNzU0Il0sInN1YiI6ImJkYmQwN2MwLWIwOWUtNDEwMi04NWQ1LTc3YWJmZmFhODg4MSIsImFtciI6WyJpYm1pZCJdLCJpYXQiOjE2Njc0MDE0MzYsInRlbmFudCI6ImE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3MiIsInNjb3BlIjoib3BlbmlkIGFwcGlkX2RlZmF1bHQgYXBwaWRfcmVhZHVzZXJhdHRyIGFwcGlkX3JlYWRwcm9maWxlIGFwcGlkX3dyaXRldXNlcmF0dHIgYXBwaWRfYXV0aGVudGljYXRlZCBzdXBlcl9lZGl0IGVkaXQgcmVhZCJ9.cFWkt1_8KtH-4DVfZw9D2gyJv6Bhagw7H3WFXZDIhZvGtH_y6t6VjJCddqnOuboA3ByxBuYreHNA5Kv6d3tD-MHZBLiBsVqkB6qBJ_SXzfRpUeXCvxnm8-O-eVVaazL0VUdyNll_NYuHqxJtrRUTPE-GMhJjmX067mCEv-iiNqEe0v42AQSEwTmBbl2Bfsp7XnYrcT00DohzYv-cbWGMYb2D0KMcCFAdDnile5T5Gz7FRcBzOP7_ZXNSbGdy2Y1JwozD2g2T-OVHW8DLOykkb5h2YdV39Mf1PyLprit5xBAyIZK6Cia9qe1L6WNp8z4QNJrmEVeL6_5bP6K04PyazA eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcElkLWE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3Mi0yMDIxLTAzLTIzVDE0OjQyOjU2LjE5NiIsInZlciI6NH0.eyJpc3MiOiJodHRwczovL2V1LWRlLmFwcGlkLmNsb3VkLmlibS5jb20vb2F1dGgvdjQvYThiZjE3YmMtMTBmNS00NzZiLWE0YzgtYjllYjFlNWQ2MDcyIiwiYXVkIjpbIjNhOWMxY2U0LTI5ODAtNGM3MS05MTcxLTJhM2NjYzQzYzc1NCJdLCJleHAiOjE2Njc0MTIyMzYsInRlbmFudCI6ImE4YmYxN2JjLTEwZjUtNDc2Yi1hNGM4LWI5ZWIxZTVkNjA3MiIsImlhdCI6MTY2NzQwMTQzNiwiZW1haWwiOiJuYXRoYW4uZ29vZEBpYm0uY29tIiwibmFtZSI6Ik5hdGhhbiBHb29kIiwic3ViIjoiYmRiZDA3YzAtYjA5ZS00MTAyLTg1ZDUtNzdhYmZmYWE4ODgxIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiTmF0aGFuLkdvb2RAaWJtLmNvbSIsImdpdmVuX25hbWUiOiJOYXRoYW4iLCJmYW1pbHlfbmFtZSI6Ikdvb2QiLCJpZGVudGl0aWVzIjpbeyJwcm92aWRlciI6ImlibWlkIiwiaWQiOiJuYXRoYW4uZ29vZEBpYm0uY29tIn1dLCJhbXIiOlsiaWJtaWQiXX0.Zz2gmW0NXMPCxZ9CmLKFV0gmxdeCKt8SPKdXwg03x2c98v1mm0eWTfF2hw6APPuZZ8BVdJI8yiwd7hlAxWpIk6dwtBC-c1UkyG23kX4iNEyVnW8XXZqpp-BZ7j88frGWdLb_R-5tsnh0GACofh6Fbp8koSphmN5OPunPRriaZAGZnESJvdiuQ63dmTXC3nU_LbaCROUyzxok2o4ohEoHRSA7-Pf9YPQOsBz5IyFJHdOZreCoE1CRqzKlsd6KLeKXQ2XbsidhTzn_oBoM-hK55dQEvr4rtyUTn1WpRRaFu7kqTKd0Xt85UAAA7Hmkpe_Wv0FqV65KN_xFJl3VBTuF7w`)
}

func TestGetOptionsAsStringsUsingHandler(t *testing.T) {
	expected := []string{
		"No",
		"Yes",
	}
	p, err := prompt.NewPromptBuilder().
		Path("root").
		Text("this is a test").
		WithOptions(prompt.YesNo()).
		Build()
	assert.NoError(t, err)
	assert.Equal(t, expected, p.OptionsToStrings())
}

func TestGetDefault(t *testing.T) {
	p, err := prompt.NewPromptBuilder().
		Path("root").
		Text("this is a test").
		AddOption("Not Default", "abc123").
		AddDefaultOption("Default", "default-option").
		Build()

	actual, exists := p.DefaultOption()
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "default-option", actual.Value())
}

func TestGetDefault_WithNone(t *testing.T) {
	p, err := prompt.NewPromptBuilder().
		Path("root").
		Text("this is a test").
		AddOption("Not Default", "abc123").
		AddOption("Also Not Default", "def456").
		Build()

	actual, exists := p.DefaultOption()
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Nil(t, actual)
}

func TestResolveEnvParams(t *testing.T) {
	type args struct {
		envvar string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test basic env var",
			args: args{
				envvar: "this-is-my-var",
			},
			want: "ITZ_THIS_IS_MY_VAR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pkg.ToEnvVar(pkg.DefaultPrefix, tt.args.envvar)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseParamDescription() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPromptsToString(t *testing.T) {
	noOptionsNoDefault, err := prompt.NewPromptBuilder().Path("test").Text("Would you like to pass this test?").Build()
	assert.NoError(t, err)

	optionsNoDefault, err := prompt.NewPromptBuilder().Path("test").Text("Would you like to pass this test?").
		WithOptions(prompt.YesNo()).
		Build()
	assert.NoError(t, err)

	optionsWithDefault, err := prompt.NewPromptBuilder().Path("test").Text("What is your favourite colour?").
		AddOption("Red", "red").
		AddDefaultOption("Blue", "blue").
		Build()
	assert.NoError(t, err)

	textOnlyWithDefault, err := prompt.NewPromptBuilder().Path("test").Text("What is your quest?").
		WithDefault("to seek the Holy Grail").
		Build()
	assert.NoError(t, err)

	type args struct {
		prompt *prompt.Prompt
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "prompt with no default and no options",
			args: args{
				prompt: noOptionsNoDefault,
			},
			want: "Would you like to pass this test?",
		},
		{
			name: "prompt with options but no default",
			args: args{
				prompt: optionsNoDefault,
			},
			want: "Would you like to pass this test? [\"No\"/\"Yes\"]",
		},
		{
			name: "prompt with no default and no options",
			args: args{
				prompt: optionsWithDefault,
			},
			want: "What is your favourite colour? [\"Blue\"/\"Red\"] (default: \"Blue\")",
		},
		{
			name: "prompt text only and a text default",
			args: args{
				prompt: textOnlyWithDefault,
			},
			want: "What is your quest? (default: \"to seek the Holy Grail\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.prompt.String()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseParamDescription() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultOptionFromHandler(t *testing.T) {
	builder := prompt.NewPromptBuilder().Path("test").Text("What option do you like?").
		WithOptions(func() (map[string]string, error) {
			results := make(map[string]string, 0)
			results["value1"] = "text1"
			results["value2"] = "text2"
			results["value3"] = "text3"
			results["value4"] = "text4"
			return results, nil
		}).
		AddDefaultOption("This is the default", "default")
	p, err := builder.Build()
	assert.NoError(t, err)
	defaultOpt, found := p.DefaultOption()
	assert.True(t, found)
	assert.Equal(t, "default", defaultOpt.Value())
	defaultVal, found := p.DefaultValue()
	assert.True(t, found)
	assert.Equal(t, "default", defaultVal)
}
