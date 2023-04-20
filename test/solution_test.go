package test

import (
	"fmt"
	"testing"

	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/stretchr/testify/assert"
)

func TestEmptyFilter(t *testing.T) {
	// test with a nil filter
	f1 := solutions.NewFilter(
		solutions.OwnerFilter(nil),
		solutions.KindFilter(nil),
	)
	assert.Empty(t, f1.BuildFilter())
	// test without calling the filter option
	f2 := solutions.NewFilter()
	assert.Empty(t, f2.BuildFilter())
}

func TestOwnerFilter(t *testing.T) {
	maximo := "ibm/ibm-maximo"
	redhat := "redhat/redhat-ansible"
	tz     := "ibm/ibm-technology-zone"
	owner  := []string{maximo, redhat, tz}
	ownerFilter := solutions.NewFilter(
		solutions.OwnerFilter(owner),
	)
	filter := ownerFilter.BuildFilter()
	expectedValue := []string{fmt.Sprintf("spec.owner=group:%s,spec.owner=group:%s,spec.owner=group:%s", maximo, redhat, tz)}
	assert.Equal(t, expectedValue, filter)
}

func TestKindFilter(t *testing.T) {
	asset := "Asset"
	colletion := "Collection"
	product := "Product"
	kind := []string{asset, colletion, product}
	kindFilter := solutions.NewFilter(
		solutions.KindFilter(kind),
	)
	filter := kindFilter.BuildFilter()
	expectedValue := []string{fmt.Sprintf("kind=%s,kind=%s,kind=%s", asset, colletion, product)}
	assert.Equal(t, expectedValue, filter)
}