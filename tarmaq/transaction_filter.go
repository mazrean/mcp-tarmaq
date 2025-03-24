package tarmaq

type TxFilter interface {
	Filter(transactions []*Transaction, query *Query) []*Transaction
}

var _ TxFilter = &MaxSizeTxFilter{}

type MaxSizeTxFilter struct {
	MaxSize int
}

func NewMaxSizeTxFilter(maxSize int) *MaxSizeTxFilter {
	return &MaxSizeTxFilter{
		MaxSize: maxSize,
	}
}

func (f *MaxSizeTxFilter) Filter(transactions []*Transaction, _ *Query) []*Transaction {
	filtered := make([]*Transaction, 0, len(transactions))

	for _, tx := range transactions {
		if len(tx.Files) <= f.MaxSize {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}

var _ TxFilter = &TarmaqTxFilter{}

//nolint:revive
type TarmaqTxFilter struct{}

func NewTarmaqTxFilter() *TarmaqTxFilter {
	return &TarmaqTxFilter{}
}

func (f *TarmaqTxFilter) Filter(transactions []*Transaction, query *Query) []*Transaction {
	k := 0
	filtered := make([]*Transaction, 0, len(transactions))

	for _, tx := range transactions {
		intersection, _ := query.Apply(tx)
		switch {
		case intersection.Len() == 0:
			continue
		case intersection.Len() == k:
			filtered = append(filtered, tx)
		case intersection.Len() > k:
			filtered = append(filtered[:0], tx)
			k = intersection.Len()
		}
	}

	return filtered
}
