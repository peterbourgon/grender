package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestDiffPath(t *testing.T) {
	type tuple struct{ base, complete string }
	for tu, expected := range map[tuple]string{
		tuple{"/foo/bar", "/foo/bar"}:             "",
		tuple{"/U/f/src", "/U/f/src/bar.html"}:    "bar.html",
		tuple{"/foo/src/", "/foo/src/a/b.json"}:   "a/b.json",
		tuple{"//foo/src/", "/foo/src/a/b.json"}:  "a/b.json",
		tuple{"/foo//src/", "/foo/src/a/b.json"}:  "a/b.json",
		tuple{"/foo/src//", "/foo/src/a/b.json"}:  "a/b.json",
		tuple{"/foo/src///", "/foo/src/a/b.json"}: "a/b.json",
	} {
		if got := Relative(tu.base, tu.complete); expected != got {
			t.Errorf("Relative(%s, %s): expected %s, got %s", tu.base, tu.complete, expected, got)
		}
	}
}

func TestMustCopy(t *testing.T) {
	src, err := ioutil.TempFile(os.TempDir(), "grender-test-mustcopy")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(src.Name())

	srcBuf := "the contents of\nthe file\n"
	if n, err := src.Write([]byte(srcBuf)); err != nil {
		t.Fatal(err)
	} else if n < len(srcBuf) {
		t.Fatalf("short write")
	}

	dst := src.Name() + ".copy"
	Copy(dst, src.Name())

	dstBuf, err := ioutil.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}

	if len(dstBuf) != len(srcBuf) {
		t.Fatalf("dst (%d) != src (%d)", len(dstBuf), len(srcBuf))
	}
	for i := 0; i < len(srcBuf); i++ {
		if dstBuf[i] != srcBuf[i] {
			t.Fatalf("dst[%d] (%d) != src[%d] (%d)", i, dstBuf[i], i, srcBuf[i])
		}
	}
}

func TestMustJSON(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "grender-test-mustjson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	buf := []byte(`{"a":"X","b":123}`)
	if err := ioutil.WriteFile(tmpFile.Name(), buf, 0655); err != nil {
		t.Fatal(err)
	}

	m := ParseJSON(Read(tmpFile.Name()))
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

func TestTargetFor(t *testing.T) {
	type tuple struct{ relativePath, ext string }
	for src, expected := range map[tuple]string{
		tuple{"/foo", ""}:            *targetDir + "/foo",
		tuple{"/foo", ".html"}:       *targetDir + "/foo.html",
		tuple{"/foo.blah", ".html"}:  *targetDir + "/foo.html",
		tuple{"/foo.html", ".blah"}:  *targetDir + "/foo.blah",
		tuple{"/a/b/c", ".php"}:      *targetDir + "/a/b/c.php",
		tuple{"/a/b/c.php", ".html"}: *targetDir + "/a/b/c.html",
	} {
		path, ext := *sourceDir+src.relativePath, src.ext
		got := TargetFor(path, ext)
		if expected != got {
			t.Errorf("%s: expected '%s', got '%s'", path, expected, got)
		}
	}
}

func TestMustTemplate(t *testing.T) {
	// TODO
	// rather a lot of setup involved here
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
		"":                 []string{}, // special-case
		"foo":              []string{"foo"},
		"/foo":             []string{"foo"},
		"foo/":             []string{"foo"},
		"a/b/c.d":          []string{"a", "b", "c.d"},
		"/foo/bar/baz.txt": []string{"foo", "bar", "baz.txt"},
	} {
		assert(SplitPath(path), expected)
	}
}

func TestDefaultTitle(t *testing.T) {
	for path, expected := range map[string]string{
		"2013-01-01-foo-bar-baz.md":    "Foo bar baz",
		"/foo/2013-1-2-foo_bar-baz.md": "Foo bar baz",
		"a/2013-a-b-foo-bar-baz":       "",
	} {
		if got := DefaultTitle(path); expected != got {
			t.Errorf("'%s': expected '%s', got '%s'", path, expected, got)
		}
	}
}

func TestDefaultDate(t *testing.T) {
	for path, expected := range map[string]string{
		"2013-01-01-foo-bar-baz.md":    "2013 01 01",
		"/foo/2013-1-2-foo_bar-baz.md": "2013 01 02",
		"a/2013-a-b-foo-bar-baz":       "",
	} {
		if got := DefaultDate(path); expected != got {
			t.Errorf("'%s': expected '%s', got '%s'", path, expected, got)
		}
	}
}

func TestSplatInto(t *testing.T) {
	m := map[string]interface{}{}
	assert := func(expected string) {
		got, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare([]byte(expected), got) != 0 {
			t.Errorf("expected '%s', got '%s'", expected, string(got))
		}
	}

	SplatInto(m, "foo", map[string]interface{}{"a": 1, "b": 2})
	assert(`{"foo":{"a":1,"b":2}}`)

	SplatInto(m, "bar/baz", map[string]interface{}{"x": map[string]string{"y": "z"}})
	assert(`{"bar":{"baz":{"x":{"y":"z"}}},"foo":{"a":1,"b":2}}`)

	SplatInto(m, "bar/baz", map[string]interface{}{"x": map[string]string{"y": "!"}})
	assert(`{"bar":{"baz":{"x":{"y":"!"}}},"foo":{"a":1,"b":2}}`)

	SplatInto(m, "bar/baz", map[string]interface{}{"x": map[string]string{"yy": "!!"}})
	assert(`{"bar":{"baz":{"x":{"y":"!","yy":"!!"}}},"foo":{"a":1,"b":2}}`)

	SplatInto(m, "foo", map[string]interface{}{"a": "x"})
	assert(`{"bar":{"baz":{"x":{"y":"!","yy":"!!"}}},"foo":{"a":"x","b":2}}`)
}
