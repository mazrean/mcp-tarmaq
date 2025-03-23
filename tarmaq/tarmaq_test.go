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

func TestTarmaq_createQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		paths   []FilePath
		fileMap map[FileID]FilePath
		wantIDs []FileID
	}{
		{
			name:    "Empty paths list",
			paths:   []FilePath{},
			fileMap: map[FileID]FilePath{},
			wantIDs: []FileID{},
		},
		{
			name: "Single path that exists in the file map",
			paths: []FilePath{
				NewFilePath("file1.txt"),
			},
			fileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
			},
			wantIDs: []FileID{1},
		},
		{
			name: "Single path that doesn't exist in the file map",
			paths: []FilePath{
				NewFilePath("nonexistent.txt"),
			},
			fileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
			},
			wantIDs: []FileID{},
		},
		{
			name: "Multiple paths that exist in the file map",
			paths: []FilePath{
				NewFilePath("file1.txt"),
				NewFilePath("file2.txt"),
			},
			fileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
				FileID(2): NewFilePath("file2.txt"),
				FileID(3): NewFilePath("file3.txt"),
			},
			wantIDs: []FileID{1, 2},
		},
		{
			name: "Mix of existing and non-existing paths",
			paths: []FilePath{
				NewFilePath("file1.txt"),
				NewFilePath("nonexistent.txt"),
				NewFilePath("file3.txt"),
			},
			fileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
				FileID(2): NewFilePath("file2.txt"),
				FileID(3): NewFilePath("file3.txt"),
			},
			wantIDs: []FileID{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tarmaq := &Tarmaq{}
			got := tarmaq.createQuery(tt.paths, tt.fileMap)

			// Check that each expected ID is in the set
			for _, id := range tt.wantIDs {
				assert.True(t, got.Files.Contains(id),
					"Expected Files set to contain ID %d", id)
			}

			// Check that the set doesn't contain any unexpected IDs from the file map
			for id := range tt.fileMap {
				contains := false
				for _, wantID := range tt.wantIDs {
					if id == wantID {
						contains = true
						break
					}
				}

				if !contains {
					assert.False(t, got.Files.Contains(id),
						"Files set should not contain ID %d", id)
				}
			}

			// Make sure the set size matches expected
			assert.Equal(t, len(tt.wantIDs), got.Files.Len(),
				"Files set has wrong size")
		})
	}
}
