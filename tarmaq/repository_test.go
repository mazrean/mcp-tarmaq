package tarmaq

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a mock Git repository
func createMockRepo(commits []mockCommit) (*git.Repository, error) {
	fs := memfs.New()
	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return nil, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	for _, commit := range commits {
		for path, content := range commit.files {
			err := func() error {
				f, err := fs.Create(path)
				if err != nil {
					return fmt.Errorf("failed to create file: %w", err)
				}
				defer f.Close()

				_, err = f.Write([]byte(content))
				if err != nil {
					return fmt.Errorf("failed to write content: %w", err)
				}

				return nil
			}()
			if err != nil {
				return nil, err
			}

			_, err = wt.Add(path)
			if err != nil {
				return nil, err
			}
		}

		_, err = wt.Commit(commit.message, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Test User",
				Email: "test@example.com",
				When:  time.Now(),
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return repo, nil
}

type mockCommit struct {
	message string
	files   map[string]string // Map of file paths and contents
}

func TestGitRepository_GetTransactions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		commits     []mockCommit
		wantTrans   []*Transaction
		wantFileMap map[FileID]FilePath
		wantErr     bool
	}{
		{
			name: "Add one file in one commit",
			commits: []mockCommit{
				{
					message: "Add file1.txt",
					files: map[string]string{
						"file1.txt": "content1",
					},
				},
			},
			wantTrans: []*Transaction{
				{
					Files: makeFileSet(FileID(0)),
				},
			},
			wantFileMap: map[FileID]FilePath{
				FileID(0): NewFilePath("file1.txt"),
			},
			wantErr: false,
		},
		{
			name: "Add different files in two commits",
			commits: []mockCommit{
				{
					message: "Add file1.txt",
					files: map[string]string{
						"file1.txt": "content1",
					},
				},
				{
					message: "Add file2.txt",
					files: map[string]string{
						"file1.txt": "content1",
						"file2.txt": "content2",
					},
				},
			},
			wantTrans: []*Transaction{
				{
					Files: makeFileSet(FileID(0)),
				},
				{
					Files: makeFileSet(FileID(1)),
				},
			},
			wantFileMap: map[FileID]FilePath{
				FileID(1): NewFilePath("file1.txt"),
				FileID(0): NewFilePath("file2.txt"),
			},
			wantErr: false,
		},
		{
			name: "Change file content",
			commits: []mockCommit{
				{
					message: "Add file1.txt",
					files: map[string]string{
						"file1.txt": "content1",
					},
				},
				{
					message: "Update file1.txt",
					files: map[string]string{
						"file1.txt": "updated content",
					},
				},
			},
			wantTrans: []*Transaction{
				{
					Files: makeFileSet(FileID(0)),
				},
				{
					Files: makeFileSet(FileID(0)),
				},
			},
			wantFileMap: map[FileID]FilePath{
				FileID(0): NewFilePath("file1.txt"),
			},
			wantErr: false,
		},
		{
			name: "Change multiple files simultaneously",
			commits: []mockCommit{
				{
					message: "Add initial files",
					files: map[string]string{
						"file1.txt": "content1",
						"file2.txt": "content2",
					},
				},
				{
					message: "Update both files",
					files: map[string]string{
						"file1.txt": "updated content1",
						"file2.txt": "updated content2",
					},
				},
			},
			wantTrans: []*Transaction{
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
				{
					Files: makeFileSet(FileID(0), FileID(1)),
				},
			},
			wantFileMap: map[FileID]FilePath{
				FileID(0): NewFilePath("file1.txt"),
				FileID(1): NewFilePath("file2.txt"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock Git repository
			repo, err := createMockRepo(tt.commits)
			if err != nil {
				t.Fatalf("failed to create mock repo: %v", err)
			}

			// Create test target object
			r := &GitRepository{
				repo: repo,
			}

			// Execute test
			gotTrans, gotFileMap, err := r.GetTransactions()

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check number of transactions
			assert.Equal(t, len(tt.wantTrans), len(gotTrans), "Number of transactions does not match")

			// Compare each transaction
			for i := 0; i < len(tt.wantTrans); i++ {
				assertSetEqual(t, tt.wantTrans[i].Files, gotTrans[i].Files, "Files of transaction %d", i)
			}

			// Compare file map
			assert.Equal(t, tt.wantFileMap, gotFileMap, "File map does not match")
		})
	}
}
