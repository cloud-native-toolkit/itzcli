package test

import (
	"os"
	"testing"
	"github.com/cloud-native-toolkit/itzcli/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestReadTechZoneUser(t *testing.T) {
	jsoner := auth.NewJsonReader()
	path, err := getPath("examples/TechZoneResponse.json")
	assert.NoError(t, err)
	fileR, err := os.Open(path)
	assert.NoError(t, err)
	user, err := jsoner.Read(fileR)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Token, "thisisarandomstring")
	assert.Equal(t, user.Preferredfirstname, "John")
	assert.Equal(t, user.Preferredlastname, "Smith")
}
