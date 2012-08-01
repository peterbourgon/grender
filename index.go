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

func (idx Index) Render() map[string][]map[string]string {
	m := map[string][]map[string]string{}
	for typ, a := range idx {
		m[typ] = make([]map[string]string, len(a))
		for i := 0; i < len(a); i++ {
			m[typ][i] = a[i].Render()
		}
	}
	return m
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
