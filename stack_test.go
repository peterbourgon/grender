package main

import (
	"github.com/peterbourgon/mergemap"
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

func TestAddForGlobal(t *testing.T) {
	assert := func(ctx string, m map[string]interface{}, key, expected string) {
		if got := m[key]; expected != got {
			t.Fatalf("%s: m['%s'] expected '%s' got '%s'", ctx, key, expected, got)
		}
	}

	s := NewStack()

	func() {
		s.Add("", map[string]interface{}{"a": map[string]interface{}{"first": "OK"}})
		v, ok := s.Get("/some/arbitrary/path")["a"]
		if !ok {
			t.Fatalf("didn't get 'a'")
		}
		m, ok := v.(map[string]interface{})
		if !ok {
			t.Fatalf("bad type for 'a'")
		}
		assert("take 1", m, "first", "OK")
	}()

	func() {
		s.Add("", map[string]interface{}{"a": map[string]interface{}{"second": "K"}})
		v, ok := s.Get("/some/other/deeper/path.html")["a"]
		if !ok {
			t.Fatalf("didn't get 'a'")
		}
		m, ok := v.(map[string]interface{})
		if !ok {
			t.Fatalf("bad type for 'a'")
		}
		assert("add 'second'", m, "first", "OK")
		assert("add 'second'", m, "second", "K")
	}()

	func() {
		s.Add("", map[string]interface{}{"a": map[string]string{"first": "NO"}})
		v, ok := s.Get("/path.md")["a"]
		if !ok {
			t.Fatalf("didn't get 'a'")
		}
		m, ok := v.(map[string]interface{})
		if !ok {
			t.Fatalf("bad type for 'a'")
		}
		assert("overwrite 'first'", m, "first", "NO")
		assert("overwrite 'first'", m, "second", "K")
	}()
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

	m1 := mergemap.Merge(m, map[string]interface{}{"a": "b"})
	if m1["a"] != "b" {
		t.Fatal("m1[a] != b")
	}

	m2 := mergemap.Merge(m1, map[string]interface{}{"a": "c"})
	if m2["a"] != "c" {
		t.Fatal("m2[a] != c")
	}

	m3 := mergemap.Merge(m2, map[string]interface{}{"b": "d"})
	if m3["a"] != "c" {
		t.Fatal("m3[a] != c")
	}
	if m3["b"] != "d" {
		t.Fatal("m3[b] != d")
	}
}
