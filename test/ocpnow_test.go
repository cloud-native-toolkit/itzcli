package test

import (
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetEnv(t *testing.T) {

	// Set up an example object with some values
	credInfo := &pkg.CredInfo{
		Name:   "helloworldcred",
		ApiKey: "thisisanapikeyexpected",
	}

	envals, err := pkg.ResolveVars(credInfo, nil)

	assert.NoError(t, err)
	assert.Equal(t, "thisisanapikeyexpected", envals["TF_VAR_ibmcloud_api_key"])
}

func TestSetEnvWithOpts(t *testing.T) {

	// Set up an example object with some values
	credInfo := &pkg.CredInfo{
		Name:   "helloworldcred",
		ApiKey: "thisisanapikeyexpected",
	}

	envals, err := pkg.ResolveVars(credInfo, &pkg.ResolveVarConfig{
		Prefix: "ATK_TF_VAR_",
	})

	assert.NoError(t, err)
	assert.Equal(t, "thisisanapikeyexpected", envals["ATK_TF_VAR_ibmcloud_api_key"])
}

func TestFindClusterByName(t *testing.T) {
	path, err := getPath("examples/exampleProject.yaml")
	assert.NoError(t, err)
	project, err := pkg.LoadProject(path)
	assert.NoError(t, err)
	cRef, err := pkg.FindClusterByName(project, "my-cluster-1")
	assert.NoError(t, err)
	assert.Equal(t, *cRef, "5fbc7b03-a094-43d8-bdd7-260d8abecf89")
	cRef, err = pkg.FindClusterByName(project, "my-cluster-2")
	assert.NoError(t, err)
	assert.Equal(t, *cRef, "288c05a4-b0eb-4d51-b3b1-676a4f1c0e19")
}

func TestNewBuildParamResolver(t *testing.T) {
	path, err := getPath("examples/exampleProject.yaml")
	assert.NoError(t, err)
	project, err := pkg.LoadProject(path)
	assert.NoError(t, err)
	params := []pkg.JobParam{
		{Name: "TF_VAR_region", Value: ""},
	}
	resolver, err := pkg.NewBuildParamResolver(project, "my-cluster-1", params)
	assert.NoError(t, err)

	actualResolved := resolver.ResolvedParams()
	assert.Equal(t, "us-east", actualResolved["TF_VAR_region"])
}
