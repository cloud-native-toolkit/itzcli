package test

import (
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGetITZDirs(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err, "failed to get user home directory")
	t.Setenv("ITZ_ERROR_DIR", filepath.Join("/tmp", "errordoesnotexist"))
	t.Setenv("ITZ_NO_ERROR_ENV_DIR", filepath.Join("/tmp"))
	tests := []struct {
		name    string
		f       func() (string, error)
		want    string
		wantErr bool
	}{
		{
			name:    "Get cache directory",
			f:       pkg.GetITZCacheDir,
			want:    filepath.Join(homeDir, ".itz", "cache"),
			wantErr: false,
		},
		{
			name:    "Get work directory",
			f:       pkg.GetITZWorkDir,
			want:    filepath.Join(homeDir, ".itz", "workspace"),
			wantErr: false,
		},
		{
			name: "Get custom directory",
			f: func() (string, error) {
				return pkg.GetITZDirOrDefault("custom", "ITZ_SOME_CUSTOM_DIR")
			},
			want:    filepath.Join(homeDir, ".itz", "custom"),
			wantErr: false,
		},
		{
			name: "Get envvar not directory with error",
			f: func() (string, error) {
				return pkg.GetITZDirOrDefault("custom", "ITZ_ERROR_DIR")
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Get envvar directory",
			f: func() (string, error) {
				return pkg.GetITZDirOrDefault("custom", "ITZ_NO_ERROR_ENV_DIR")
			},
			want:    "/tmp",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f()
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

func TestKeyify(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "hyphens should be removed",
			args: args{
				name: "ocp-installer",
			},
			want: "ocpinstaller",
		},
		{
			name: "should be all lowercase",
			args: args{
				name: "ocpInstaller",
			},
			want: "ocpinstaller",
		},
		{
			name: "underscores should be removed",
			args: args{
				name: "ocp_installer",
			},
			want: "ocpinstaller",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pkg.Keyify(tt.args.name); got != tt.want {
				t.Errorf("Keyify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenCommandName(t *testing.T) {
	type args struct {
		cmd    *cobra.Command
		suffix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple command",
			args: args{
				cmd:    createSimpleCmd(),
				suffix: "default",
			},
			want: "parent.child.default",
		},
		{
			name: "should be all lowercase",
			args: args{
				cmd:    createSimpleCmd(),
				suffix: "install",
			},
			want: "parent.child.install",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pkg.FlattenCommandName(tt.args.cmd, tt.args.suffix); got != tt.want {
				t.Errorf("Keyify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createSimpleCmd() *cobra.Command {
	subcmd := &cobra.Command{
		Use: "child",
	}
	cmd := &cobra.Command{
		Use: "parent",
	}
	cmd.AddCommand(subcmd)
	return subcmd
}
