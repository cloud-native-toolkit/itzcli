package test

import (
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"testing"
)

func TestAppendToFilename(t *testing.T) {
	type args struct {
		fn     string
		suffix string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal happy path replacement",
			args: args{
				fn:     "https://github.com/cloud-native-toolkit/deployer-operator-masauto/blob/main/maximo-pipeline.yaml",
				suffix: "-run",
			},
			want:    "https://github.com/cloud-native-toolkit/deployer-operator-masauto/blob/main/maximo-pipeline-run.yaml",
			wantErr: false,
		},
		{
			name: "normal happy path replacement - file scheme",
			args: args{
				fn:     "file:///home/user1/cloud-native-toolkit/deployer-operator-masauto/blob/main/maximo-pipeline.yaml",
				suffix: "-run",
			},
			want:    "file:///home/user1/cloud-native-toolkit/deployer-operator-masauto/blob/main/maximo-pipeline-run.yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pkg.AppendToFilename(tt.args.fn, tt.args.suffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppendToFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AppendToFilename() got = %v, want %v", got, tt.want)
			}
		})
	}
}
