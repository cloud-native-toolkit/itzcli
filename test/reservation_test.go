package test

import (
	"bytes"
	"reflect"

	"os"
	"path/filepath"
	"testing"

	reservations "github.com/cloud-native-toolkit/itzcli/pkg/techzone"
	"github.com/stretchr/testify/assert"
)

type fileReservationServiceClient struct {
	FilePath string
}

func (c *fileReservationServiceClient) Get(id string) (*reservations.Reservation, error) {
	return nil, nil
}

func (c *fileReservationServiceClient) GetAll(f reservations.Filter) ([]reservations.Reservation, error) {
	jsoner := reservations.NewJsonReader()
	fileR, err := os.Open(c.FilePath)
	rez, err := jsoner.ReadAll(fileR)
	// Now filter them
	result := make([]reservations.Reservation, 0)
	if f != nil {
		for _, r := range rez {
			if f(r) {
				result = append(result, r)
			}
		}
	}
	return result, err
}

func TestReadReservations(t *testing.T) {
	path, err := getPath("examples/myReservationsResponse.json")
	client := &fileReservationServiceClient{
		FilePath: path,
	}

	rez, err := client.GetAll(reservations.NoFilter())
	assert.NoError(t, err)
	assert.NotNil(t, rez)
	assert.Equal(t, len(rez), 4)
}

func TestFilterReadyReservations(t *testing.T) {
	path, err := getPath("examples/myReservationsResponse.json")
	client := &fileReservationServiceClient{
		FilePath: path,
	}

	rez, err := client.GetAll(reservations.FilterByStatus("Deleted"))

	assert.NoError(t, err)
	assert.NotNil(t, rez)
	assert.Equal(t, 1, len(rez))

	tw := reservations.NewModelWriter(reflect.TypeOf(reservations.Reservation{}).Name(), "text")

	// TODO: Use a mock client and a filter
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	err = tw.WriteMany(buf, rez)
	assert.NoError(t, err)

	assert.Equal(t,
		" - Request IBM zSystems access - Deleted\n   Reservation Id: 5e4QKpDt3l96T1f5Lz7PYhSR\n\n",
		buf.String())

}

func TestJSONFormat(t *testing.T) {
	path, err := getPath("examples/myReservationsResponse.json")
	client := &fileReservationServiceClient{
		FilePath: path,
	}

	rez, err := client.GetAll(reservations.FilterByStatusSlice([]string{"Deleted", "Expired"}))

	assert.NoError(t, err)
	assert.NotNil(t, rez)
	assert.Equal(t, 2, len(rez))

	tw := reservations.NewModelWriter(reflect.TypeOf(reservations.Reservation{}).Name(), "json")

	// TODO: Use a mock client and a filter
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	err = tw.WriteMany(buf, rez)
	assert.NoError(t, err)
	jsonString := `[{"CollectionId":"5j5DlbGrwd9bOKTdo0o59YbC","CreatedAt":1678199733811,"Description":"Example of asking for access","ExtendCount":0,"Name":"Request IBM zSystems access","ProvisionDate":"2023-03-07 14:35:00","ProvisionUntil":"2023-03-11 14:35:00","id":"5e4QKpDt3l96T1f5Lz7PYhSR","ServiceLinks":[{"type":"desktop","Label":"Desktop client","Sensitive":false,"Url":"https://desktop.download.ibm.com/client"},{"type":"password","Label":"Desktop Client Password","Sensitive":true,"Url":"thesecretpassword"},{"type":"username","Label":"Desktop Client Username","Sensitive":false,"Url":"myuser"},{"type":"public_endpoint","Label":"Public Client Endpoint","Sensitive":false,"Url":"192.168.10.10:8088"}],"Status":"Deleted"},{"CollectionId":"5j5DlbGrwd9bOKTdo0o59YbC","CreatedAt":1678199733811,"Description":"Example of asking for access","ExtendCount":0,"Name":"Request IBM zSystems access","ProvisionDate":"2023-03-07 14:35:00","ProvisionUntil":"2023-03-11 14:35:00","id":"5e4QKpDt3l96T1f5Lz7PYhSR","ServiceLinks":[{"type":"desktop","Label":"Desktop client","Sensitive":false,"Url":"https://desktop.download.ibm.com/client"},{"type":"password","Label":"Desktop Client Password","Sensitive":true,"Url":"thesecretpassword"},{"type":"username","Label":"Desktop Client Username","Sensitive":false,"Url":"myuser"},{"type":"public_endpoint","Label":"Public Client Endpoint","Sensitive":false,"Url":"192.168.10.10:8088"}],"Status":"Expired"}]`

	assert.Equal(t,
		jsonString,
		buf.String())
}

func getPath(name string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(workingDir, name), nil
}
