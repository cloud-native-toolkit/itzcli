package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloud-native-toolkit/itzcli/api"
	"github.com/stretchr/testify/assert"

)

func TestLogin(t *testing.T) {
	r := api.SetUpRouter()
	missingTokenreq, err := http.NewRequest("GET", "/login", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, missingTokenreq)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

    badAccessTokenreq, err := http.NewRequest("GET", "/login?token=12345", nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, badAccessTokenreq)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
