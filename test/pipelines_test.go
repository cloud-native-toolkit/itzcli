package test

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/mocks"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"reflect"
	"testing"
)

func TestIsDefault(t *testing.T) {

	tests := []struct {
		name string
		args *pkg.PipelineParamOption
		want bool
	}{
		{
			name: "Test true",
			args: &pkg.PipelineParamOption{
				Default: "true",
			},
			want: true,
		},
		{
			name: "Test not populated",
			args: &pkg.PipelineParamOption{},
			want: false,
		},
		{
			name: "Test false",
			args: &pkg.PipelineParamOption{
				Default: "false",
			},
			want: false,
		},
		{
			name: "Test false",
			args: &pkg.PipelineParamOption{
				Default: "waitwotthisisnotavalidboolean",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.IsDefault())
		})
	}
}

func TestParseParamDescription(t *testing.T) {
	type args struct {
		from string
	}
	tests := []struct {
		name    string
		args    args
		want    *pkg.PipelineParamParts
		wantErr bool
	}{
		{
			name: "Test three options with one default",
			args: args{
				from: `specify the preferred storageclass
{
	"options": [
	{"text": "thin","value": "thin", "default": "true"},
	{"text": "gp2","value": "gp2" },
	{"text": "ocs-storagecluster-cephfs","value": "ocs-storagecluster-cephfs" }
	]
}	
`,
			},
			want: &pkg.PipelineParamParts{
				Description: "specify the preferred storageclass",
				ParamOptions: pkg.PipelineParamOptions{
					Options: []pkg.PipelineParamOption{
						{
							Text:    "thin",
							Value:   "thin",
							Default: "true",
						},
						{
							Text:  "gp2",
							Value: "gp2",
						},
						{
							Text:  "ocs-storagecluster-cephfs",
							Value: "ocs-storagecluster-cephfs",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test description without any JSON",
			args: args{
				"The IBM entitlement key with permissions for pulling cp4d images\n",
			},
			want: &pkg.PipelineParamParts{
				Description: "The IBM entitlement key with permissions for pulling cp4d images",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pkg.ParseParamDescription(tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseParamDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseParamDescription() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildPipelinePrompt_AcceptDefaults(t *testing.T) {
	// Load up the Pipeline from a file
	client := &pkg.GitServiceClient{
		BaseDest: "/tmp",
	}
	path, err := getPath("examples/deployerPipeline.yaml")
	assert.NoError(t, err, "expected no error getting a path")
	gitRepo := fmt.Sprintf("file://%s", path)
	t.Log(fmt.Sprintf("Using %s for the deployer file path...", gitRepo))
	pipeline, err := client.Get(gitRepo, pkg.UnmarshalPipeline)
	pl := pipeline.(*v1beta1.Pipeline)
	assert.NoError(t, err, "expected no error getting a pipeline")
	assert.NotNil(t, pipeline)

	pipelineResolver := pkg.NewPipelineResolver(pl)

	// Then pass it to the parser to get the Prompt
	prompt, err := pkg.BuildPipelinePrompt(pl.Name, pipelineResolver, pipelineResolver)
	assert.NoError(t, err, "expected no error creating a prompt from the pipeline")
	assert.NotNil(t, prompt)
	assert.Equal(t, 1, len(prompt.SubPrompts()), "there should only be one parameter required when using defaults")
}

func TestBuildPipelinePrompt_AcceptDefaultsWithEnvVar(t *testing.T) {
	// Load up the Pipeline from a file
	client := &pkg.GitServiceClient{
		BaseDest: "/tmp",
	}
	path, err := getPath("examples/deployerPipeline.yaml")
	assert.NoError(t, err, "expected no error getting a path")
	gitRepo := fmt.Sprintf("file://%s", path)
	t.Log(fmt.Sprintf("Using %s for the deployer file path...", gitRepo))
	pipeline, err := client.Get(gitRepo, pkg.UnmarshalPipeline)
	pl := pipeline.(*v1beta1.Pipeline)
	assert.NoError(t, err, "expected no error getting a pipeline")
	assert.NotNil(t, pipeline)

	pipelineResolver := pkg.NewPipelineResolver(pl)
	envResolver := pkg.NewEnvParamResolver()
	chainedResolver := pkg.NewChainedResolver(pkg.UseEnvironmentVars|pkg.UsePipelineDefaults, envResolver, pipelineResolver)

	// Now set up the environment with ITZ_NAMESPACE
	t.Setenv("ITZ_NAMESPACE", "my-test-pipeline")

	// Then pass it to the parser to get the Prompt
	prompt, err := pkg.BuildPipelinePrompt(pl.Name, pipelineResolver, chainedResolver)
	assert.NoError(t, err, "expected no error creating a prompt from the pipeline")
	assert.NotNil(t, prompt)
	assert.Equal(t, 0, len(prompt.SubPrompts()), "there should be no parameter required when using defaults and env set")
}

func TestBuildPipelinePrompt(t *testing.T) {
	// Load up the Pipeline from a file
	client := &pkg.GitServiceClient{
		BaseDest: "/tmp",
	}
	path, err := getPath("examples/deployerPipeline.yaml")
	assert.NoError(t, err, "expected no error getting a path")
	gitRepo := fmt.Sprintf("file://%s", path)
	t.Log(fmt.Sprintf("Using %s for the deployer file path...", gitRepo))
	pipeline, err := client.Get(gitRepo, pkg.UnmarshalPipeline)
	pl := pipeline.(*v1beta1.Pipeline)
	assert.NoError(t, err, "expected no error getting a pipeline")
	assert.NotNil(t, pipeline)

	pipelineResolver := pkg.NewPipelineResolver(pl)

	// Then pass it to the parser to get the Prompt
	prompt, err := pkg.BuildPipelinePrompt(pl.Name, pipelineResolver, pkg.NewEnvParamResolver())
	assert.NoError(t, err, "expected no error creating a prompt from the pipeline")
	assert.NotNil(t, prompt)
	assert.Equal(t, 65, len(prompt.SubPrompts()), "there should be 65 parameters when not using the default resolver")
}

func TestArgResolver(t *testing.T) {
	argParser := pkg.NewArgsParamParser([]string{"namespace=fred", "malformedohno", "variable-1=value-1", "variable-2=oh=my", "variable-3=\"with quotes and spaces\""})
	assert.True(t, argParser.EnabledFor(pkg.UseCommandLineArgs))
	assert.False(t, argParser.EnabledFor(pkg.UsePipelineDefaults))
	val, exists := argParser.Lookup("namespace")
	assert.True(t, exists)
	assert.Equal(t, "fred", val)

	tests := []struct {
		name        string
		param       string
		want        string
		shouldExist bool
	}{
		{
			name:        "namespace variable should exist and be set to \"fred\"",
			param:       "namespace",
			want:        "fred",
			shouldExist: true,
		},
		{
			name:        "testing non-existing arg",
			param:       "doh",
			want:        "",
			shouldExist: false,
		},
		{
			name:        "testing second variable after error for robustness",
			param:       "variable-1",
			want:        "value-1",
			shouldExist: true,
		},
		{
			name:        "testing variable with equal sign in actual value",
			param:       "variable-2",
			want:        "oh=my",
			shouldExist: true,
		},
		{
			name:        "testing variable with spaces",
			param:       "variable-3",
			want:        "\"with quotes and spaces\"",
			shouldExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, exists := argParser.Lookup(tt.param)
			if exists != tt.shouldExist {
				t.Errorf("Lookup() exists = %v, shouldExist %v", exists, tt.shouldExist)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lookup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainedResolver(t *testing.T) {
	mockEnvResolver := mocks.NewParamResolver(t)
	mockParamResolver := mocks.NewParamResolver(t)
	options := pkg.UseEnvironmentVars

	mockEnvResolver.On("EnabledFor", pkg.UseEnvironmentVars).Return(true)
	mockEnvResolver.On("EnabledFor", pkg.UsePipelineDefaults).Return(false)
	mockParamResolver.On("EnabledFor", pkg.UsePipelineDefaults).Return(false)
	chainedResolver := pkg.NewChainedResolver(options, mockEnvResolver, mockParamResolver)
	assert.True(t, chainedResolver.EnabledFor(pkg.UseEnvironmentVars))
	assert.False(t, chainedResolver.EnabledFor(pkg.UsePipelineDefaults))

	expected := "my-value"
	mockEnvResolver.On("EnabledFor", pkg.UseEnvironmentVars).Return(true)
	mockEnvResolver.On("Lookup", "my-variable").Return(expected, true)
	actual, exists := chainedResolver.Lookup("my-variable")
	assert.True(t, exists)
	assert.Equal(t, expected, actual)
	mockParamResolver.AssertNotCalled(t, "Lookup", "this option is not supported, so should not be called.")

	mockEnvResolver.On("EnabledFor", pkg.UseEnvironmentVars).Return(true)
	mockEnvResolver.On("Lookup", "my-variable-2").Return("", false)
	mockParamResolver.On("EnabledFor", pkg.UseEnvironmentVars).Return(false)
	actual, exists = chainedResolver.Lookup("my-variable-2")
	assert.False(t, exists)
	assert.Equal(t, "", actual)

}

func TestPipelineRunMarshall(t *testing.T) {
	// Load up the Pipeline from a file
	client := &pkg.GitServiceClient{
		BaseDest: "/tmp",
	}
	path, err := getPath("examples/examplePipelineRun.yaml")
	assert.NoError(t, err, "expected no error getting a path")
	gitRepo := fmt.Sprintf("file://%s", path)
	t.Log(fmt.Sprintf("Using %s for the deployer file path...", gitRepo))
	pipeline, err := client.Get(gitRepo, pkg.UnmarshalPipelineRun)
	prun := pipeline.(*v1beta1.PipelineRun)
	assert.NoError(t, err, "expected no error getting a pipeline")
	assert.NotNil(t, prun)
	assert.True(t, pkg.IsPipelineRun(*prun))
	assert.Equal(t, 4, len(prun.Spec.Params))
}

func TestPipelineRunMerge(t *testing.T) {
	mockReader := mocks.NewParamReader(t)
	mockResolver := mocks.NewParamResolver(t)

	// Set up the mocks.
	mockReader.On("Params").Return([]v1beta1.ParamSpec{
		// new
		{
			Name: "my-new-param-1",
			Default: &v1beta1.ParamValue{
				Type:      v1beta1.ParamTypeString,
				StringVal: "my-default-param-value-1",
			},
		},
		// new
		{
			Name: "my-new-param-2",
			Default: &v1beta1.ParamValue{
				Type:      v1beta1.ParamTypeString,
				StringVal: "my-default-param-value-2",
			},
		},
		// updated
		{
			Name: "tag-name",
		},
		// updated
		{
			Name: "repo-url",
		},
	}, nil)

	mockResolver.On("Lookup", "my-new-param-1").Return("my-param-value-1", true)
	mockResolver.On("Lookup", "my-new-param-2").Return("my-param-value-2", true)
	mockResolver.On("Lookup", "tag-name").Return("develop", true)
	mockResolver.On("Lookup", "repo-url").Return("https://github.com/cloud-native-toolkit/itzcli", true)

	// Load up the Pipeline from a file
	client := &pkg.GitServiceClient{
		BaseDest: "/tmp",
	}
	path, err := getPath("examples/examplePipelineRun.yaml")
	assert.NoError(t, err, "expected no error getting a path")
	gitRepo := fmt.Sprintf("file://%s", path)
	t.Log(fmt.Sprintf("Using %s for the deployer file path...", gitRepo))
	pipeline, err := client.Get(gitRepo, pkg.UnmarshalPipelineRun)
	prun := pipeline.(*v1beta1.PipelineRun)
	assert.NoError(t, err, "expected no error getting a pipeline")
	assert.NotNil(t, prun)
	assert.True(t, pkg.IsPipelineRun(*prun))
	assert.Equal(t, 4, len(prun.Spec.Params))

	pl := &v1beta1.Pipeline{}

	merged, err := pkg.MergePipelineRun(prun, pl, mockReader, mockResolver)
	assert.NoError(t, err, "expected no error getting a pipeline")
	assert.NotNil(t, merged)

	// The merged pipeline run should have now 6 params--two of them new, and two
	// of them updated from the original values. The other two should be the same
	// as the original values.
	assert.Equal(t, 6, len(merged.Spec.Params))
	updatedParam, found := pkg.FindParam(merged.Spec.Params, "tag-name")
	assert.True(t, found)
	assert.Equal(t, "develop", updatedParam.Value.StringVal)
}
