package tarmaq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTarmaq_createResults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		rules       []*Rule
		fileMap     map[FileID]FilePath
		wantResults []*Result
	}{
		{
			name:        "Empty rules list",
			rules:       []*Rule{},
			fileMap:     map[FileID]FilePath{},
			wantResults: []*Result{},
		},
		{
			name: "Single rule with valid file ID",
			rules: []*Rule{
				{
					Right:      FileID(1),
					Confidence: 0.8,
					Support:    10,
				},
			},
			fileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
			},
			wantResults: []*Result{
				{
					Path:       NewFilePath("file1.txt"),
					Confidence: 0.8,
					Support:    10,
				},
			},
		},
		{
			name: "Multiple rules with the same target file",
			rules: []*Rule{
				{
					Right:      FileID(1),
					Confidence: 0.8,
					Support:    10,
				},
				{
					Right:      FileID(1),
					Confidence: 0.9,
					Support:    15,
				},
			},
			fileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
			},
			wantResults: []*Result{
				{
					Path:       NewFilePath("file1.txt"),
					Confidence: 0.8,
					Support:    10,
				},
			},
		},
		{
			name: "Rule with invalid file ID",
			rules: []*Rule{
				{
					Right:      FileID(99),
					Confidence: 0.8,
					Support:    10,
				},
			},
			fileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
			},
			wantResults: []*Result{},
		},
		{
			name: "Rules with mix of valid and invalid file IDs",
			rules: []*Rule{
				{
					Right:      FileID(1),
					Confidence: 0.8,
					Support:    10,
				},
				{
					Right:      FileID(99),
					Confidence: 0.9,
					Support:    15,
				},
				{
					Right:      FileID(2),
					Confidence: 0.7,
					Support:    5,
				},
			},
			fileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
				FileID(2): NewFilePath("file2.txt"),
			},
			wantResults: []*Result{
				{
					Path:       NewFilePath("file1.txt"),
					Confidence: 0.8,
					Support:    10,
				},
				{
					Path:       NewFilePath("file2.txt"),
					Confidence: 0.7,
					Support:    5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tarmaq := &Tarmaq{}
			gotResults := tarmaq.createResults(tt.rules, tt.fileMap)

			// Map iteration order is non-deterministic, so check element matching
			assert.ElementsMatch(t, tt.wantResults, gotResults)
		})
	}
}
