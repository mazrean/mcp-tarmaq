package tarmaq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxSizeTxFilter_Filter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		maxSize      int
		transactions []*Transaction
		query        *Query
		want         []*Transaction
	}{
		{
			name:    "All transactions below max size",
			maxSize: 3,
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1), FileID(2)),
				},
			},
			query: &Query{
				Files: makeFileSet(),
			},
			want: []*Transaction{
				{
					Files: makeFileSet(FileID(0)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1), FileID(2)),
				},
			},
		},
		{
			name:    "Some transactions above max size",
			maxSize: 2,
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1), FileID(2)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1), FileID(2), FileID(3)),
				},
			},
			query: &Query{
				Files: makeFileSet(),
			},
			want: []*Transaction{
				{
					Files: makeFileSet(FileID(0)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
			},
		},
		{
			name:    "All transactions above max size",
			maxSize: 1,
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1), FileID(2)),
				},
			},
			query: &Query{
				Files: makeFileSet(),
			},
			want: []*Transaction{},
		},
		{
			name:    "Edge case - max size zero",
			maxSize: 0,
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0)),
				},
				{
					Files: makeFileSet(),
				},
			},
			query: &Query{
				Files: makeFileSet(),
			},
			want: []*Transaction{
				{
					Files: makeFileSet(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewMaxSizeTxFilter(tt.maxSize)
			got := filter.Filter(tt.transactions, tt.query)

			assert.Equal(t, len(tt.want), len(got), "Number of filtered transactions does not match")

			// Compare each transaction
			for i := 0; i < len(tt.want); i++ {
				if i >= len(got) {
					break
				}
				assertSetEqual(t, tt.want[i].Files, got[i].Files, "Files of transaction %d", i)
			}
		})
	}
}

func TestTarmaqTxFilter_Filter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		transactions []*Transaction
		query        *Query
		want         []*Transaction
	}{
		{
			name: "No matching transactions",
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(2), FileID(3)),
				},
			},
			query: &Query{
				Files: makeFileSet(FileID(4), FileID(5)),
			},
			want: []*Transaction{},
		},
		{
			name: "Single matching transaction",
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(2)),
				},
				{
					Files: makeFileSet(FileID(3), FileID(4)),
				},
			},
			query: &Query{
				Files: makeFileSet(FileID(1)),
			},
			want: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(2)),
				},
			},
		},
		{
			name: "Different intersection sizes - keep max",
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(2), FileID(3)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(3), FileID(4)),
				},
			},
			query: &Query{
				Files: makeFileSet(FileID(1), FileID(3)),
			},
			want: []*Transaction{
				{
					Files: makeFileSet(FileID(1), FileID(2), FileID(3)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(3), FileID(4)),
				},
			},
		},
		{
			name: "Empty query",
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(2), FileID(3)),
				},
			},
			query: &Query{
				Files: makeFileSet(),
			},
			want: []*Transaction{},
		},
		{
			name: "Multiple transactions with same max intersection",
			transactions: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(2)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(3)),
				},
			},
			query: &Query{
				Files: makeFileSet(FileID(1)),
			},
			want: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(2)),
				},
				{
					Files: makeFileSet(FileID(1), FileID(3)),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewTarmaqTxFilter()
			got := filter.Filter(tt.transactions, tt.query)

			assert.Equal(t, len(tt.want), len(got), "Number of filtered transactions does not match")

			// Compare each transaction
			for i := 0; i < len(tt.want); i++ {
				if i >= len(got) {
					break
				}
				assertSetEqual(t, tt.want[i].Files, got[i].Files, "Files of transaction %d", i)
			}
		})
	}
}
