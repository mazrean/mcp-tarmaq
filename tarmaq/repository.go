package tarmaq

import (
	"errors"
	"fmt"
	"log/slog"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mazrean/mcp-tarmaq/pkg/collection"
)

type Repository interface {
	GetTransactions() ([]*Transaction, map[FileID]FilePath, error)
}

var _ Repository = &GitRepository{}

type GitRepository struct {
	repo             *git.Repository
	transactionLimit int
}

func NewGitRepository(repoPath string, transactionLimit int) (*GitRepository, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	return &GitRepository{
		repo:             repo,
		transactionLimit: transactionLimit,
	}, nil
}

func (r *GitRepository) GetTransactions() ([]*Transaction, map[FileID]FilePath, error) {
	commitIter, err := r.repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("get commit iterator: %w", err)
	}
	defer commitIter.Close()

	idGenerator := FileIDGenerator{0}
	var transactions []*Transaction
	latestFileMap := make(map[FileID]FilePath)
	fileIDMap := make(map[string]FileID)

	for commit, err := commitIter.Next(); err == nil; commit, err = commitIter.Next() {
		files := collection.NewSet[FileID]()

		var parentTree *object.Tree
		// get first parent(main branch in most cases)
		parent, err := commit.Parent(0)
		switch {
		case errors.Is(err, object.ErrParentNotFound):
			// empty tree for first commit
			parentTree = &object.Tree{}
		case err == nil:
			parentTree, err = parent.Tree()
			if err != nil {
				slog.Warn("failed to get parent tree",
					slog.String("commit", commit.Hash.String()),
					slog.String("parent", parent.Hash.String()),
					slog.String("error", err.Error()),
				)
				continue
			}
		default:
			slog.Warn("failed to get parent",
				slog.String("commit", commit.Hash.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		commitTree, err := commit.Tree()
		if err != nil {
			slog.Warn("failed to get commit tree",
				slog.String("commit", commit.Hash.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		changes, err := parentTree.Diff(commitTree)
		if err != nil {
			slog.Warn("failed to get diff",
				slog.String("commit", commit.Hash.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		for _, change := range changes {
			// add file to transaction if it's added or modified
			fileID, ok := fileIDMap[change.To.Name]
			if !ok {
				fileID = idGenerator.Next()
				latestFileMap[fileID] = NewFilePath(change.To.Name)
				fileIDMap[change.To.Name] = fileID
			}
			files.Add(fileID)

			// edit fileIDMap if file is renamed
			if change.From.Name != change.To.Name {
				delete(fileIDMap, change.To.Name)
				fileIDMap[change.From.Name] = fileID
			}
		}

		if files.Len() > 0 {
			transactions = append(transactions, &Transaction{
				Files: files,
			})
			if r.transactionLimit != 0 && len(transactions) >= r.transactionLimit {
				break
			}
		}
	}

	return transactions, latestFileMap, nil
}
