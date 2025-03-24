package tarmaq

import (
	"slices"
	"testing"

	"github.com/mazrean/mcp-tarmaq/pkg/collection"
	"github.com/stretchr/testify/assert"
)

// Helper function for Set
func makeFileSet(ids ...FileID) collection.Set[FileID] {
	s := collection.NewSet[FileID]()
	for _, id := range ids {
		s.Add(id)
	}
	return s
}

// Function to convert Set to slice and assert
func assertSetEqual(t *testing.T, expected, actual collection.Set[FileID], msgAndArgs ...any) {
	expectedSlice := slices.Collect(expected.Iter())
	actualSlice := slices.Collect(actual.Iter())
	assert.ElementsMatch(t, expectedSlice, actualSlice, msgAndArgs...)
}
