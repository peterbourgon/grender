package main

import (
	"testing"
)

func TestAddGet(t *testing.T) {
	s := NewStack()
	s.Add("/a", map[string]interface{}{"a": "A"})
	s.Add("/a/1", map[string]interface{}{"one": "1"})
	s.Add("/b", map[string]interface{}{"b": "B"})
	t.Logf("%v", s.m)

	assert := func(m map[string]interface{}, key, expected string) {
		v, ok := m[key]
		if !ok {
			t.Fatalf("key '%s' wasn't present", key)
		}
		s, ok := v.(string)
		if !ok {
			t.Fatalf("key '%s' wasn't a string", key)
		}
		if s != expected {
			t.Fatalf("key '%s' wasn't '%s'", key, expected)
		}
	}

	assert(s.Get("/a/foo.txt"), "a", "A")
}

func TestSplitPath(t *testing.T) {
	assert := func(a, b []string) {
		if len(a) != len(b) {
			t.Fatalf("%v != %v", a, b)
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				t.Fatalf("%v != %v", a, b)
			}
		}
	}

	for path, expected := range map[string][]string{
		"":                 []string{},
		"foo":              []string{"foo"},
		"/foo":             []string{"foo"},
		"foo/":             []string{"foo"},
		"a/b/c.d":          []string{"a", "b", "c.d"},
		"/foo/bar/baz.txt": []string{"foo", "bar", "baz.txt"},
	} {
		assert(splitPath(path), expected)
	}
}

func TestMergeInto(t *testing.T) {
	m := map[string]interface{}{}

	m1 := mergeInto(m, map[string]interface{}{"a": "b"})
	if m1["a"] != "b" {
		t.Fatal("m1[a] != b")
	}

	m2 := mergeInto(m1, map[string]interface{}{"a": "c"})
	if m2["a"] != "c" {
		t.Fatal("m2[a] != c")
	}

	m3 := mergeInto(m2, map[string]interface{}{"b": "d"})
	if m3["a"] != "c" {
		t.Fatal("m3[a] != c")
	}
	if m3["b"] != "d" {
		t.Fatal("m3[b] != d")
	}
}
