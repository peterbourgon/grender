package main

// An Index is a map of types to ordered index-tuples.
type Index map[string]OrderedIndexTuples

func (idx Index) Append(typ string, tuple *IndexTuple) {
	if a, ok := idx[typ]; ok {
		idx[typ] = append(a, tuple)
	} else {
		idx[typ] = []*IndexTuple{tuple}
	}
}

//
//
//

// OrderedIndexTuples is a sortable list of IndexTuples.
// It's sorted by the *indexSortKey.
type OrderedIndexTuples []*IndexTuple

func (a OrderedIndexTuples) Len() int { return len(a) }

func (a OrderedIndexTuples) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a OrderedIndexTuples) Less(i, j int) bool {
	return a[i].SortKey > a[j].SortKey
}
