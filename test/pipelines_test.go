package test

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/stretchr/testify/assert"
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

func TestBuildPipelinePrompt(t *testing.T) {
	// Load up the Pipeline from a file
	client := &pkg.GitServiceClient{
		BaseDest: "/tmp",
	}
	path, err := getPath("examples/deployerPipeline.yaml")
	assert.NoError(t, err, "expected no error getting a path")
	gitRepo := fmt.Sprintf("file://%s", path)
	t.Log(fmt.Sprintf("Using %s for the deployer file path...", gitRepo))
	pipeline, err := client.Get(gitRepo)
	assert.NoError(t, err, "expected no error getting a pipeline")
	assert.NotNil(t, pipeline)

	// Then pass it to the parser to get the Prompt
	prompt, err := pkg.BuildPipelinePrompt(pipeline)
	assert.NoError(t, err, "expected no error creating a prompt from the pipeline")
	assert.NotNil(t, prompt)
	assert.Equal(t, 65, len(prompt.SubPrompts()))
}
