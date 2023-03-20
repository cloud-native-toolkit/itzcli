package test

import (
	"fmt"
	"testing"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/stretchr/testify/assert"
)



func TestEmptyQueryParam(t *testing.T) {
	// test with a nil query
	q1 := solutions.NewQuery(
		solutions.OwnerQuery(nil),
	)
	assert.Empty(t, q1.BuildQuery())
	// test without calling the query option  
	q2 := solutions.NewQuery()
	assert.Empty(t, q2.BuildQuery())
}


func TestQueryParam(t *testing.T) {
	commandLine := []string{"partners/palantir", "partners/mongodb"}
	query := solutions.NewQuery(
		solutions.OwnerQuery(commandLine),
	)
	fmt.Print(query.BuildQuery())
	assert.Equal(t, "filter=spec.owner=group:partners/palantir,partners/mongodb", query.BuildQuery())
}