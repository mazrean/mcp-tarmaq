package tarmaq

import (
	"iter"

	"github.com/mazrean/mcp-tarmaq/pkg/collection"
)

type Extractor interface {
	Extract(transactions []*Transaction, query *Query) []*Rule
}

var _ Extractor = &AssociationRuleExtractor{}

type AssociationRuleExtractor struct {
	minConfidence float64
	minSupport    uint64
}

func NewAssociationRuleExtractor(minConfidence float64, minSupport uint64) *AssociationRuleExtractor {
	return &AssociationRuleExtractor{
		minConfidence: minConfidence,
		minSupport:    minSupport,
	}
}

func (e *AssociationRuleExtractor) Extract(transactions []*Transaction, query *Query) []*Rule {
	supportMap := make(SupportMap)
	for _, tx := range transactions {
		left, rights := query.Apply(tx)
		if left.Len() == 0 {
			continue
		}

		supportMapItem := supportMap.Load(left)
		for right := range rights.Iter() {
			supportMapItem.ruleMap[right]++
		}
		supportMapItem.support++
	}

	rules := []*Rule{}
	for rule := range supportMap.Iter() {
		if rule.Confidence >= e.minConfidence && rule.Support >= e.minSupport {
			rules = append(rules, rule)
		}
	}

	return rules
}

type SupportMap map[uint64][]*SupportMapItem

type SupportMapItem struct {
	left    collection.Set[FileID]
	ruleMap map[FileID]uint64
	support uint64
}

func (s SupportMap) Load(left collection.Set[FileID]) *SupportMapItem {
	items, ok := s[left.Hash()]
	if !ok {
		items = make([]*SupportMapItem, 0)
		s[left.Hash()] = items
	}

	for _, item := range items {
		if item.left.Equal(left) {
			return item
		}
	}

	item := &SupportMapItem{
		left:    left,
		ruleMap: make(map[FileID]uint64),
		support: 0,
	}
	items = append(items, item)
	s[left.Hash()] = items

	return item
}

func (s SupportMap) Iter() iter.Seq[*Rule] {
	return func(yield func(*Rule) bool) {
		for _, items := range s {
			for _, item := range items {
				for right, support := range item.ruleMap {
					confidence := float64(support) / float64(item.support)
					if !yield(&Rule{
						Left:       item.left,
						Right:      right,
						Confidence: confidence,
						Support:    support,
					}) {
						return
					}
				}
			}
		}
	}
}
