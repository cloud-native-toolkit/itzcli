package test

import (
	"github.com/spf13/cobra"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"testing"
)

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
