package test

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/cloud-native-toolkit/itzcli/cmd/dr"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfigCheck_DoCheck(t *testing.T) {
	type fields struct {
		ConfigKey string
		Defaulter dr.DefaultGetter
	}
	type args struct {
		tryFix bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &dr.ConfigCheck{
				ConfigKey: tt.fields.ConfigKey,
				Defaulter: tt.fields.Defaulter,
			}
			got, err := c.DoCheck(tt.args.tryFix)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DoCheck() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoChecks(t *testing.T) {
	type args struct {
		checks []dr.Check
		tryFix bool
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dr.DoChecks(tt.args.checks, tt.args.tryFix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileCheck_DoCheck(t *testing.T) {
	type fields struct {
		Path  string
		Name  string
		IsDir bool
		Help  string
	}
	type args struct {
		tryFix bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &dr.FileCheck{
				Path:  tt.fields.Path,
				Name:  tt.fields.Name,
				IsDir: tt.fields.IsDir,
				Help:  tt.fields.Help,
			}
			got, err := f.DoCheck(tt.args.tryFix)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DoCheck() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetITZHomeDir(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dr.GetITZHomeDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetITZHomeDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetITZHomeDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceFileCheck(t *testing.T) {
	type args struct {
		c 	 dr.CheckerFunc
		help string
		f    dr.FileAutoFixFunc
	}
	real_c := dr.OneExistsOnPath("podman","docker")
	real_f := dr.UpdateConfig("podman.path")
	tests := []struct {
		name string
		args args
		want *dr.FileCheck
	}{
		{	
			name: 		"Test if creates FileCheck object",
			args: args{
				c: 		real_c,
				help:   "%s was not found on your path",
				f: 		real_f,
			},
			want: &dr.FileCheck{
				PathCheckFunc: 	real_c,
				Path:           os.Getenv("PATH"),
				Name:        	"",
				IsDir:       	false,
				Help:        	"%s was not found on your path",
				UpdaterFunc: 	real_f,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dr.NewResourceFileCheck(tt.args.c, tt.args.help, tt.args.f);
			v 	:= reflect.ValueOf(got).Interface().(*dr.FileCheck)
			wv	:= reflect.ValueOf(tt.want).Interface().(*dr.FileCheck)
			
			assert.Equal(t, v.Help, wv.Help, "Help is not equal")
			assert.Equal(t, v.Path, wv.Path, "Path is not equal")
			assert.Equal(t, v.IsDir, wv.IsDir, "IsDir is not equal")
			assert.Equal(t, v.Name, wv.Name, "Name is not equal")
		})
	}
}

func TestNewConfigCheck(t *testing.T) {
	type args struct {
		configKey string
		help      string
		defaulter dr.DefaultGetter
	}
	tests := []struct {
		name string
		args args
		want *dr.ConfigCheck
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dr.NewConfigCheck(tt.args.configKey, "", tt.args.defaulter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfigDirCheck(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *dr.FileCheck
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dr.NewReqConfigDirCheck(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigDirCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfigFileCheck(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *dr.FileCheck
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dr.NewConfigFileCheck(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigFileCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoDefault(t *testing.T) {
	tests := []struct {
		name string
		want dr.DefaultGetter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dr.NoDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceURL(t *testing.T) {
	type args struct {
		scheme string
		port   int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "http test to 8080",
			args: args{
				scheme: "http",
				port:   8080,
			},
		},
		{
			name: "https test to 8088",
			args: args{
				scheme: "https",
				port:   8088,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// It is really hard to compare these, so we will just make sure
			// that it is a valid URL
			actualUrl := dr.ServiceURL(tt.args.scheme, tt.args.port)().(string)
			parsedUrl, err := url.Parse(actualUrl)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.scheme, parsedUrl.Scheme)
			assert.Equal(t, fmt.Sprintf("%d", tt.args.port), strings.Split(parsedUrl.Host, ":")[1])
		})
	}
}

func TestStatic(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want dr.DefaultGetter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dr.Static(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Static() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIPAddressForPodmanConnection(t *testing.T) {
	getter := dr.ServiceURL("http", 8080)
	wdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	viper.Set("podman.path", filepath.Join(wdir, "scripts/mock_podman"))
	t.Log(viper.Get("podman.path"))
	value := getter().(string)
	assert.Equal(t, "http://172.16.16.128:8080", value)
}
