package main

import (
	"fmt"
	"path/filepath"
)

type Stack struct {
	m map[string]map[string]interface{} // path: partial-metadata
}

func NewStack() *Stack {
	return &Stack{
		m: map[string]map[string]interface{}{},
	}
}

func (s *Stack) Add(path string, m map[string]interface{}) {
	println("Add", path, "yields", filepath.Clean(path))
	s.m[filepath.Clean(path)] = m
}

func (s *Stack) Get(path string) map[string]interface{} {
	m := map[string]interface{}{}
	list := splitPath(path)
	if len(list) <= 0 {
		return m
	}
	fmt.Printf("Get '%s' split to '%v'\n", path, list)

	for i, _ := range list {
		key := filepath.Join(list[:i+1]...)
		println("Get", path, "getting", key)
		if m0, ok := s.m[key]; ok {
			m = mergeInto(m, m0)
		}
	}
	return m
}

func mergeInto(dst, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
