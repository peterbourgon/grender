package main

import (
	"path/filepath"
	"strings"
)

// Stack stores a set of keyed map[string]interface{} metadata.
// The key is specified during Add, and is assumed to be a valid
// file path. Get supplies a path, and returns the union of all
// Added metadata along every step in the given path.
//
// For example, Get("/foo/bar/baz") returns merged metadata for "/foo",
// "/foo/bar", and "/foo/bar/baz". In this way, the Stack enables the
// 'stackable' Grender context behavior.
type Stack struct {
	m map[string]map[string]interface{} // path: partial-metadata
}

func NewStack() *Stack {
	return &Stack{
		m: map[string]map[string]interface{}{},
	}
}

func (s *Stack) Add(path string, m map[string]interface{}) {
	key := filepath.Join(splitPath(path)...)

	existing, ok := s.m[key]
	if !ok {
		existing = map[string]interface{}{}
	}

	s.m[key] = mergeInto(existing, m)
}

func (s *Stack) Get(path string) map[string]interface{} {
	list := splitPath(path)
	if len(list) <= 0 {
		return map[string]interface{}{}
	}

	// A weird bit of trickery. We add global keys with a 'path' of "" (empty
	// string) with the expectation that Get will return them for every input
	// path. So, we prepend "" to every lookup request. That means 'i' is off-
	// by-one, so we can use it directly against the list slice.
	m := map[string]interface{}{}
	for i, _ := range append([]string{""}, list...) {
		key := filepath.Join(list[:i]...)
		if m0, ok := s.m[key]; ok {
			m = mergeInto(m, m0)
		}
	}
	return m
}

// splitPath tokenizes the given path string on filepath.Separator.
func splitPath(path string) []string {
	list := []string{}
	for _, s := range strings.Split(path, string(filepath.Separator)) {
		if s := strings.TrimSpace(s); s != "" {
			list = append(list, s)
		}
	}
	return list
}

// mergeInto merges the src map into the tgt map, returning the union.
// Key collisions are handled by preferring src.
func mergeInto(tgt, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// TODO recursive merge of maps to some reasonable depth
		tgt[k] = v
	}
	return tgt
}
