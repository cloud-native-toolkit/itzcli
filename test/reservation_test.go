package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloud-native-toolkit/itzcli/pkg/reservations"
	"github.com/stretchr/testify/assert"
)

func TestReadReservations(t *testing.T) {
	jsoner := reservations.NewJsonReader()
	path, err := getPath("examples/reservationsResponse.json")
	assert.NoError(t, err)
	fileR, err := os.Open(path)
	assert.NoError(t, err)
	rez, err := jsoner.ReadAll(fileR)
	assert.NoError(t, err)
	assert.NotNil(t, rez)
	assert.Equal(t, len(rez), 7)
}

func TestFilterReadyReservations(t *testing.T) {
	jsoner := reservations.NewJsonReader()
	path, err := getPath("examples/reservationsResponse.json")
	assert.NoError(t, err)
	fileR, err := os.Open(path)
	assert.NoError(t, err)
	rez, err := jsoner.ReadAll(fileR)

	assert.NoError(t, err)
	assert.NotNil(t, rez)
	assert.Equal(t, len(rez), 7)

	tw := reservations.NewTextWriter()

	// HACK: I wanted to use mock testing here, but the mock is really hard
	// to set up with the template because it's not as straightforward as
	// just counting the number of times that io.SolutionWriter.Write() was called,
	// which I was hoping was the case.
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	match, err := tw.WriteFilter(buf, rez, reservations.FilterByStatus("Ready"))
	assert.NoError(t, err)
    if match == 0 {
        t.Errorf("Should at least have one match from WriteFilter")
    }

	// TODO: We might want to compare this against a file.
	assert.Equal(t, buf.String(), " - RedHat 8.4 Base Image (Fyre Advanced) - Ready\n   Reservation Id: 6320ab2046c677001874e1be\n\n - Redhat 8.5 Base Image with RDP (Fyre-2) - Ready\n   Reservation Id: 6320b5b346c677001874e1c4\n\n")

}

func getPath(name string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(workingDir, name), nil
}
