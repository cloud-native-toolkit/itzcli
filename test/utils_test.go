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
	"github.com/cloud-native-toolkit/itzcli/pkg"
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
			got, err := pkg.GetITZHomeDir()
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
		c    dr.CheckerFunc
		help string
		f    dr.FileAutoFixFunc
	}
	real_c := dr.OneExistsOnPath("podman", "docker")
	real_f := dr.UpdateConfig("podman.path")
	tests := []struct {
		name string
		args args
		want *dr.FileCheck
	}{
		{
			name: "Test if creates FileCheck object",
			args: args{
				c:    real_c,
				help: "%s was not found on your path",
				f:    real_f,
			},
			want: &dr.FileCheck{
				PathCheckFunc: real_c,
				Path:          os.Getenv("PATH"),
				Name:          "",
				IsDir:         false,
				Help:          "%s was not found on your path",
				UpdaterFunc:   real_f,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dr.NewResourceFileCheck(tt.args.c, tt.args.help, tt.args.f)
			v := reflect.ValueOf(got).Interface().(*dr.FileCheck)
			wv := reflect.ValueOf(tt.want).Interface().(*dr.FileCheck)

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

func TestOneExistsOnPath(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name      string
		args      args
		wantPath  string
		wantBin   string
		wantFound bool
	}{
		{
			name: "Tests basic paths that should exist",
			args: args{
				[]string{"bash", "cat", "find"},
			},
			wantPath:  "/bin",
			wantBin:   "bash",
			wantFound: true,
		},
		{
			name: "Tests exactly one that should not exist",
			args: args{
				[]string{"whatdoesthecowsay"},
			},
			wantPath:  "",
			wantBin:   "",
			wantFound: false,
		},
		{
			name: "Tests the second one should exist",
			args: args{
				[]string{"whatdoesthecowsay", "cat"},
			},
			wantPath:  "/bin",
			wantBin:   "cat",
			wantFound: true,
		},
		{
			name: "Tests should stop at the first found",
			args: args{
				[]string{"cat", "whatdoesthecowsay"},
			},
			wantPath:  "/bin",
			wantBin:   "cat",
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotBin, gotFound := dr.OneExistsOnPath(tt.args.names...)()
			if len(gotPath) == 0 && len(tt.wantPath) > 0 {
				t.Errorf("OneExistsOnPath() = %v, want %v", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotBin, tt.wantBin) {
				t.Errorf("OneExistsOnPath() = %v, want %v", gotBin, tt.wantBin)
			}
			if gotFound != tt.wantFound {
				t.Errorf("OneExistsOnPath() = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestActionDoCheck(t *testing.T) {
	type args struct {
		logMsg   string
		preCheck dr.PreChecker
		cmd      dr.ActionRunner
		tryFix   bool
	}
	tests := []struct {
		name    string
		args    args
		wantMsg string
		wantErr bool
	}{
		{
			name: "test precheck says false",
			args: args{
				logMsg: "precheck returned false",
				preCheck: func() bool {
					return false
				},
				cmd: func() (string, error) {
					return "this should not run", nil
				},
			},
			wantMsg: "skipping action as precheck is false",
			wantErr: false,
		},
		{
			name: "test precheck passes but action errors",
			args: args{
				tryFix: true,
				logMsg: "precheck passes",
				preCheck: func() bool {
					return true
				},
				cmd: func() (string, error) {
					return "action error", fmt.Errorf("there was an error in the action")
				},
			},
			wantMsg: "action error",
			wantErr: true,
		},
		{
			name: "test precheck passes and action succeeds",
			args: args{
				logMsg: "precheck passes",
				preCheck: func() bool {
					return true
				},
				cmd: func() (string, error) {
					return "", nil
				},
			},
			wantMsg: "",
			wantErr: false,
		},
		{
			name: "test no precheck",
			args: args{
				logMsg:   "no precheck",
				preCheck: nil,
				cmd: func() (string, error) {
					return "", nil
				},
			},
			wantMsg: "",
			wantErr: false,
		},
		{
			name: "test no cmd",
			args: args{
				logMsg:   "no cmd",
				preCheck: nil,
				cmd:      nil,
			},
			wantMsg: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := dr.NewCmdActionCheck("logged message", tt.args.preCheck, tt.args.cmd)
			gotMsg, gotErr := check.DoCheck(tt.args.tryFix)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMsg, tt.wantMsg) {
				t.Errorf("Check() = %v, want %v", gotMsg, tt.wantMsg)
			}

		})
	}
}
