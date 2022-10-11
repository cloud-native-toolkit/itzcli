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
