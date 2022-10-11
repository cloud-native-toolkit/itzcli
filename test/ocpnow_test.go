package test

import (
	"github.com/stretchr/testify/assert"
	"github.ibm.com/skol/atkcli/pkg"
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
	path, err := getPath("examples/project.yaml")
	assert.NoError(t, err)
	project, err := pkg.LoadProject(path)
	assert.NoError(t, err)
	cRef, err := pkg.FindClusterByName(project, "gartnerdemoibm")
	assert.NoError(t, err)
	assert.Equal(t, *cRef, "93370a12-490c-410c-9854-c2fcd11483b8")
	cRef, err = pkg.FindClusterByName(project, "myawsdemo")
	assert.NoError(t, err)
	assert.Equal(t, *cRef, "2f4a3119-11ee-403e-8710-fef0358e938c")
}
