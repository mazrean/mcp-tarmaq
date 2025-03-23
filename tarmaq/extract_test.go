package tarmaq

import (
	"testing"

	"github.com/mazrean/mcp-tarmaq/pkg/collection"
	"github.com/stretchr/testify/assert"
)

func TestAssociationRuleExtractor_Extract(t *testing.T) {
	tests := []struct {
		name          string
		minConfidence float64
		minSupport    uint64
		transactions  []*Transaction
		query         *Query
		expectedRules []*Rule
	}{
		{
			name:          "Simple test case",
			minConfidence: 0.5,
			minSupport:    2,
			transactions: []*Transaction{
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3))},
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(4))},
				{Files: collection.NewSet(FileID(1), FileID(3), FileID(5))},
				{Files: collection.NewSet(FileID(2), FileID(3), FileID(6))},
			},
			query: &Query{Files: collection.NewSet(FileID(1))},
			expectedRules: []*Rule{
				{
					Left:       collection.NewSet(FileID(1)),
					Right:      FileID(2),
					Confidence: 0.66666,
					Support:    2,
				},
				{
					Left:       collection.NewSet(FileID(1)),
					Right:      FileID(3),
					Confidence: 0.66666,
					Support:    2,
				},
				// Other rules are not included because they fall below the confidence or support threshold
			},
		},
		{
			name:          "Query with multiple files",
			minConfidence: 0.6,
			minSupport:    2,
			transactions: []*Transaction{
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3), FileID(4))},
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3), FileID(5))},
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(5), FileID(6))},
				{Files: collection.NewSet(FileID(1), FileID(3), FileID(4), FileID(7))},
			},
			query: &Query{Files: collection.NewSet(FileID(1), FileID(2))},
			expectedRules: []*Rule{
				{
					Left:       collection.NewSet(FileID(1), FileID(2)),
					Right:      FileID(3),
					Confidence: 0.66666,
					Support:    2,
				},
				{
					Left:       collection.NewSet(FileID(1), FileID(2)),
					Right:      FileID(5),
					Confidence: 0.66666,
					Support:    2,
				},
			},
		},
		{
			name:          "High threshold settings",
			minConfidence: 0.9,
			minSupport:    3,
			transactions: []*Transaction{
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3))},
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3))},
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(4))},
				{Files: collection.NewSet(FileID(1), FileID(3), FileID(5))},
			},
			query:         &Query{Files: collection.NewSet(FileID(1))},
			expectedRules: []*Rule{
				// No rules meet the high threshold conditions for confidence and support
			},
		},
		{
			name:          "No transactions matching the query",
			minConfidence: 0.5,
			minSupport:    1,
			transactions: []*Transaction{
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3))},
				{Files: collection.NewSet(FileID(2), FileID(3), FileID(4))},
				{Files: collection.NewSet(FileID(3), FileID(4), FileID(5))},
			},
			query:         &Query{Files: collection.NewSet(FileID(10))},
			expectedRules: []*Rule{
				// No rules are extracted because there are no transactions matching the query
			},
		},
		{
			name:          "Larger dataset",
			minConfidence: 0.6,
			minSupport:    3,
			transactions: []*Transaction{
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3), FileID(4))},
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3), FileID(5))},
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(3), FileID(6))},
				{Files: collection.NewSet(FileID(1), FileID(2), FileID(4), FileID(7))},
				{Files: collection.NewSet(FileID(1), FileID(3), FileID(5), FileID(8))},
				{Files: collection.NewSet(FileID(2), FileID(4), FileID(6), FileID(9))},
				{Files: collection.NewSet(FileID(1), FileID(3), FileID(5), FileID(10))},
			},
			query: &Query{Files: collection.NewSet(FileID(1))},
			expectedRules: []*Rule{
				{
					Left:       collection.NewSet(FileID(1)),
					Right:      FileID(2),
					Confidence: 0.66666,
					Support:    4,
				},
				{
					Left:       collection.NewSet(FileID(1)),
					Right:      FileID(3),
					Confidence: 0.83333,
					Support:    5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewAssociationRuleExtractor()
			rules := extractor.Extract(tt.transactions, tt.query, tt.minConfidence, tt.minSupport)

			assert.Len(t, rules, len(tt.expectedRules))

			// 各期待ルールが実際の結果に含まれることを確認
			for _, expectedRule := range tt.expectedRules {
				found := false
				for _, rule := range rules {
					if expectedRule.Left.Equal(rule.Left) &&
						expectedRule.Right == rule.Right &&
						assert.InDeltaf(t, expectedRule.Confidence, rule.Confidence, 0.001, "confidence mismatch") &&
						expectedRule.Support == rule.Support {
						found = true
						break
					}
				}
				assert.Truef(t, found, "expected rule not found: %+v", expectedRule)
			}
		})
	}
}
