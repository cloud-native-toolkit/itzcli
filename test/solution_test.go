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
	assert.Len(t, f1.BuildFilter(), 0)
	// test without calling the filter option
	f2 := solutions.NewFilter()
	assert.Len(t, f2.BuildFilter(), 0)
}

const (
	maximo    = "ibm/ibm-maximo"
	redhat    = "redhat/redhat-ansible"
	tz        = "ibm/ibm-technology-zone"
	asset     = "Asset"
	colletion = "Collection"
	product   = "Product"
)

func TestOwnerFilter(t *testing.T) {
	owner := []string{maximo, redhat, tz}
	ownerFilter := solutions.NewFilter(
		solutions.OwnerFilter(owner),
	)
	filter := ownerFilter.BuildFilter()
	expectedValue := []string{expectedOwner()}
	assert.Equal(t, expectedValue, filter)
}

func TestKindFilter(t *testing.T) {
	kind := []string{asset, colletion, product}
	kindFilter := solutions.NewFilter(
		solutions.KindFilter(kind),
	)
	filter := kindFilter.BuildFilter()
	expectedValue := []string{expectedKind()}
	assert.Equal(t, expectedValue, filter)
}

func TestAllFilter(t *testing.T) {
	owner := []string{maximo, redhat, tz}
	kind := []string{asset, colletion, product}
	filter := solutions.NewFilter(
		solutions.OwnerFilter(owner),
		solutions.KindFilter(kind),
	).BuildFilter()
	expectedValue := []string{expectedOwner() + fmt.Sprintf(",%s", expectedKind())}
	assert.Equal(t, expectedValue, filter)

	filter2 := solutions.NewFilter(
		solutions.KindFilter(kind),
		solutions.OwnerFilter(owner),
	).BuildFilter()
	expectedValue2 := []string{expectedKind() + fmt.Sprintf(",%s", expectedOwner())}
	assert.Equal(t, expectedValue2, filter2)
}

func expectedOwner() string {
	return fmt.Sprintf("spec.owner=%s,spec.owner=%s,spec.owner=%s", maximo, redhat, tz)
}

func expectedKind() string {
	return fmt.Sprintf("kind=%s,kind=%s,kind=%s", asset, colletion, product)
}
