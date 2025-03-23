package tarmaq

import (
	"path/filepath"

	"github.com/mazrean/mcp-tarmaq/pkg/collection"
)

type FileID uint64

type FileIDGenerator struct {
	nextID FileID
}

func (f *FileIDGenerator) Next() FileID {
	id := f.nextID
	f.nextID++
	return id
}

type FilePath string

func NewFilePath(path string) FilePath {
	return FilePath(filepath.FromSlash(path))
}

type Query struct {
	Files collection.Set[FileID]
}

func (q *Query) Apply(transaction *Transaction) (intersection collection.Set[FileID], difference collection.Set[FileID]) {
	return transaction.Files.Intersection(q.Files), transaction.Files.Difference(q.Files)
}

type Transaction struct {
	Files collection.Set[FileID]
}

type Rule struct {
	Left       collection.Set[FileID]
	Right      FileID
	Confidence float64
	Support    uint64
}

func (r *Rule) Apply(query *Query) bool {
	return r.Left.Subset(query.Files)
}
