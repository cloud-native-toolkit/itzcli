package test

import (
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
	commandLine := []string{"partners/palantir"}
	singleQuery := solutions.NewQuery(
		solutions.OwnerQuery(commandLine),
	)
	assert.Equal(t, "filter=spec.owner=group:partners/palantir", singleQuery.BuildQuery())
	commandLine = append(commandLine, "partners/mongodb")
	twoQuery := solutions.NewQuery(
		solutions.OwnerQuery(commandLine),
	)
	assert.Equal(t, "filter=spec.owner=group:partners/palantir,spec.owner=group:partners/mongodb", twoQuery.BuildQuery())
	commandLine = append(commandLine, "partners/cockroachdb")
	threeQuery := solutions.NewQuery(
		solutions.OwnerQuery(commandLine),
	)
	assert.Equal(t, "filter=spec.owner=group:partners/palantir,spec.owner=group:partners/mongodb,spec.owner=group:partners/cockroachdb", threeQuery.BuildQuery())
}