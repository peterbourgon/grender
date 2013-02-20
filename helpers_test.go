package main

import (
	"io/ioutil"
	"os"
	"testing"
)

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
		"/foo/bar/baz.txt": []string{"foo", "bar", "baz.txt"},
	} {
		assert(splitPath(path), expected)
	}
}

func TestDiffPath(t *testing.T) {
	type tuple struct{ base, complete string }
	for tu, expected := range map[tuple]string{
		tuple{"/U/f/src", "/U/f/src/bar.html"}:    "bar.html",
		tuple{"/foo/src/", "/foo/src/a/b.json"}:   "a/b.json",
		tuple{"//foo/src/", "/foo/src/a/b.json"}:  "a/b.json",
		tuple{"/foo//src/", "/foo/src/a/b.json"}:  "a/b.json",
		tuple{"/foo/src//", "/foo/src/a/b.json"}:  "a/b.json",
		tuple{"/foo/src///", "/foo/src/a/b.json"}: "a/b.json",
	} {
		if got := diffPath(tu.base, tu.complete); expected != got {
			t.Errorf("diffPath(%s, %s): expected %s, got %s", tu.base, tu.complete, expected, got)
		}
	}
}

func TestCopyFile(t *testing.T) {
	// TODO
}

func TestReadJSON(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "grender-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	buf := []byte(`{"a":"X","b":123}`)
	if err := ioutil.WriteFile(tmpFile.Name(), buf, 0655); err != nil {
		t.Fatal(err)
	}

	m := mustJSON(mustRead(tmpFile.Name()))
	a, ok := m["a"]
	if !ok {
		t.Fatal("'a' not present")
	}
	s, ok := a.(string)
	if !ok {
		t.Fatal("'a' not string")
	}
	if s != "X" {
		t.Fatal("'a' not 'X'")
	}

	b, ok := m["b"]
	if !ok {
		t.Fatal("'b' not present")
	}
	i, ok := b.(float64)
	if !ok {
		t.Fatal("'b' not number")
	}
	if i != 123 {
		t.Fatal("'b' not 123")
	}
}
