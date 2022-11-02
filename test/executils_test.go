package test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.ibm.com/skol/itzcli/pkg"
	"testing"
)

func TestImageFound(t *testing.T) {
	out := new(bytes.Buffer)
	out.WriteString("localhost/bifrost:latest\n")
	out.WriteString("localhost/atkci:latest\n")
	assert.True(t, pkg.ImageFound(out, "localhost/bifrost"))
}

func TestImageNotFound(t *testing.T) {
	out := new(bytes.Buffer)
	out.WriteString("localhost/bifrost:latest\n")
	out.WriteString("localhost/atkci:latest\n")
	assert.False(t, pkg.ImageFound(out, "localhost/mooshoopork"))
}

func TestImageWithQuotes(t *testing.T) {
	out := new(bytes.Buffer)
	out.WriteString("\"localhost/bifrost:latest\"\n")
	out.WriteString("\"localhost/atkci:latest\"\n")
	assert.True(t, pkg.ImageFound(out, "localhost/bifrost"))
}

func TestImageWithLatestTag(t *testing.T) {
	out := new(bytes.Buffer)
	out.WriteString("\"localhost/bifrost:latest\"\n")
	out.WriteString("\"localhost/atkci:latest\"\n")
	assert.True(t, pkg.ImageFound(out, "localhost/bifrost:latest"))
}
