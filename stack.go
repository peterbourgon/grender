package main

import (
	"path/filepath"

	"github.com/peterbourgon/mergemap"
)

type StackReader interface {
	Get(path string) map[string]interface{}
}

type StackWriter interface {
	Add(path string, m map[string]interface{})
}

type StackReadWriter interface {
	StackReader
	StackWriter
}

// Stack stores a set of keyed map[string]interface{} metadata. The key is
// specified during Add, and is assumed to be a valid file path. Get supplies a
// path, and returns the union of all Added metadata along every step in the
// given path.
//
// Metadata added with a path of "" (empty string) is considered global, and
// will be returned with every Get.
//
// As an example, Get("/foo/bar/baz") returns merged metadata for "", "/foo",
// "/foo/bar", and "/foo/bar/baz", preferring keys from more explicit (deeper)
// paths. In this way, Stack enables the 'stackable' Grender context behavior.
type Stack struct {
	m map[string]map[string]interface{} // path: partial-metadata
}

func NewStack() *Stack {
	return &Stack{
		m: map[string]map[string]interface{}{},
	}
}

// Add merges the given metadata into the Stack element represented by path.
// If no such element exists, Add will create it.
func (s *Stack) Add(path string, m map[string]interface{}) {
	key := filepath.Join(SplitPath(path)...)

	existing, ok := s.m[key]
	if !ok {
		existing = map[string]interface{}{}
	}

	s.m[key] = mergemap.Merge(existing, m)
}

// Get returns the aggregate metadata visible from the given path.
func (s *Stack) Get(path string) map[string]interface{} {
	list := SplitPath(path)
	if len(list) <= 0 {
		return map[string]interface{}{}
	}

	// A weird bit of trickery. We add global metadata with a path of "" (empty
	// string) under the expectation that Get will return them for every input
	// path. So, we prepend "" to every lookup request. That means 'i' is off-
	// by-one, so we can use it directly against the list slice.
	m := map[string]interface{}{}
	for i, _ := range append([]string{""}, list...) {
		key := filepath.Join(list[:i]...)
		if m0, ok := s.m[key]; ok {
			m = mergemap.Merge(m, m0)
		}
	}
	return m
}
