package test

import (
	"github.ibm.com/skol/itzcli/pkg"
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
