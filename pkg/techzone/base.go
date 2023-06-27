package techzone

import (
	"encoding/json"
	"io"

	"github.com/cloud-native-toolkit/itzcli/pkg"
)

type ServiceLink struct {
	LinkType  string `json:"type"`
	Label     string
	Sensitive bool
	Url       string
}

type Reservation struct {
	Name         string
	ServiceLinks []ServiceLink
	//OpportunityId  string
	ReservationId  string `json:"id"`
	CreatedAt      int
	Status         string
	ProvisionDate  string
	ProvisionUntil string
	CollectionId   string
	ExtendCount    int
	Description    string
}

type Filter func(Reservation) bool

func NoFilter() Filter {
	return func(r Reservation) bool {
		return true
	}
}

func FilterByStatus(status string) Filter {
	return func(r Reservation) bool {
		return r.Status == status
	}
}

func FilterByStatusSlice(status []string) Filter {
	return func(r Reservation) bool {
		return pkg.StringSliceContains(status, r.Status)
	}
}

type Reader interface {
	Read(io.Reader) (Reservation, error)
	ReadAll(io.Reader) ([]Reservation, error)
}

type JsonReader struct{}

func (j *JsonReader) Read(reader io.Reader) (Reservation, error) {
	var res Reservation
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func (j *JsonReader) ReadAll(reader io.Reader) ([]Reservation, error) {
	var res []Reservation
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func NewJsonReader() *JsonReader {
	return &JsonReader{}
}
