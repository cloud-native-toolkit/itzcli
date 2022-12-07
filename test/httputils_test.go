package test

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestGetPort(t *testing.T) {
	u, err := url.Parse("http://192.168.111.128:8088")
	assert.NoError(t, err)
	assert.Equal(t, "8088", u.Port())
}
