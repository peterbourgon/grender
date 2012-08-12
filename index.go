package main

import (
	"sort"
)

type Index []*SourceFile

func (idx *Index) Add(sf *SourceFile) {
	*idx = append(*idx, sf)
	sort.Sort(*idx)
}

func (idx Index) Len() int      { return len(idx) }
func (idx Index) Swap(i, j int) { idx[i], idx[j] = idx[j], idx[i] }
func (idx Index) Less(i, j int) bool {
	return idx[i].getString(*sortkeyKey) > idx[j].getString(*sortkeyKey)
}

func (idx *Index) Render(n int) []map[string]interface{} {
	m := []map[string]interface{}{}
	for _, sf := range *idx {
		m = append(m, sf.Render())
		if len(m) >= n {
			break
		}
	}
	return m
}
