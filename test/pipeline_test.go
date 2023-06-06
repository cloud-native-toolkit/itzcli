package test

import (
	"testing"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/stretchr/testify/assert"
)

// TestBuildDestination tests that the destination directory is built correctly
// from the GitHub URL.
//
// The directory should be the path after the GitHub organization and repository
// name, minus the ".git" suffix. As an example, if the GitHub URL is
// https://github.com/cloud-native-toolkit/itzcli/dir/myfile.yaml, then the
// result should be ${basedir}/dir/myfile.yaml, where ${basedir} is the value
// supplied as base in BuildDestination
func TestBuildDestination(t *testing.T) {
	gitRepo := "https://github.com/cloud-native-toolkit/deployer-cloud-pak-deployer/blob/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml"
	actual, err := pkg.BuildDestination("/home/user/.itz/cache", gitRepo)
	assert.NoError(t, err, "should be no error building destination")
	assert.Equal(t, "/home/user/.itz/cache/cloud-native-toolkit/deployer-cloud-pak-deployer/blob/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml", actual)
}

func TestMapGitUrlToRaw(t *testing.T) {
	gitRepo := "https://github.com/cloud-native-toolkit/deployer-cloud-pak-deployer/blob/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml"
	raw, err := pkg.MapGitUrlToRaw(gitRepo)
	assert.NoError(t, err, "should be no error mapping Git URL to raw")
	assert.Equal(t, "https://raw.githubusercontent.com/cloud-native-toolkit/deployer-cloud-pak-deployer/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml", raw)
}

func TestMapGitUrlToRaw_AlreadyRaw(t *testing.T) {
	gitRepo := "https://raw.githubusercontent.com/cloud-native-toolkit/deployer-cloud-pak-deployer/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml"
	raw, err := pkg.MapGitUrlToRaw(gitRepo)
	assert.NoError(t, err, "should be no error mapping Git URL to raw that is already raw")
	assert.Equal(t, "https://raw.githubusercontent.com/cloud-native-toolkit/deployer-cloud-pak-deployer/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml", raw)
}

func TestMapGitUrlToRaw_Malformed(t *testing.T) {
	gitRepo := "https://github.com/cloud-native-toolkit/deployer-cloud-pak-deployer"
	raw, err := pkg.MapGitUrlToRaw(gitRepo)
	assert.Error(t, err, "should be error mapping Git URL to raw")
	assert.Empty(t, raw)
}

func TestGetGitPipeline(t *testing.T) {
	client := &pkg.GitServiceClient{
		BaseDest: "/tmp",
	}

	gitRepo := "https://github.com/cloud-native-toolkit/deployer-cloud-pak-deployer/blob/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml"
	pipeline, err := client.Get(gitRepo)
	assert.NoError(t, err)
	assert.Equal(t, "cloud-pak-deployer", pipeline.Name())
	assert.FileExists(t, "/tmp/cloud-native-toolkit/deployer-cloud-pak-deployer/main/openshift-4.10/cp4d-4.6.4/cloud-pak-deployer.yaml")
}
