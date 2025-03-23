package tarmaq

import "log/slog"

type Tarmaq struct {
	Repository Repository
	TxFilters  []TxFilter
	Extractor  Extractor
}

func NewTarmaq(repo Repository, txFilters []TxFilter, extractor Extractor) *Tarmaq {
	return &Tarmaq{
		Repository: repo,
		TxFilters:  txFilters,
		Extractor:  extractor,
	}
}

type Result struct {
	Path       FilePath
	Confidence float64
	Support    uint64
}

func (t *Tarmaq) Execute(query *Query) ([]*Result, error) {
	transactions, fileMap, err := t.Repository.GetTransactions()
	if err != nil {
		return nil, err
	}

	for _, filter := range t.TxFilters {
		transactions = filter.Filter(transactions, query)
	}

	rules := t.Extractor.Extract(transactions, query)

	return t.createResults(rules, fileMap), nil
}

func (t *Tarmaq) createResults(rules []*Rule, fileMap map[FileID]FilePath) []*Result {
	resultMap := make(map[FileID]*Result)
	for _, rule := range rules {
		if _, ok := resultMap[rule.Right]; ok {
			continue
		}

		path, ok := fileMap[rule.Right]
		if !ok {
			slog.Warn("file not found",
				slog.Uint64("file_id", uint64(rule.Right)),
			)
			continue
		}
		resultMap[rule.Right] = &Result{
			Path:       path,
			Confidence: rule.Confidence,
			Support:    rule.Support,
		}
	}

	results := make([]*Result, 0, len(resultMap))
	for _, result := range resultMap {
		results = append(results, result)
	}

	return results
}
