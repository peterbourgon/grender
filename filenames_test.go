package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBasename(t *testing.T) {
	m := map[string]struct {
		prefix   string
		filename string
	}{
		"foo":     {"", "foo"},
		"alpha":   {"", "alpha.txt"},
		"bar/baz": {"", "bar/baz.source"},
		"beta":    {"/a/", "/a/beta"},
		"a/b/c":   {"/path/to/source", "/path/to/source/a/b/c.txt"},
	}

	for expected, pair := range m {
		got := Basename(pair.prefix, pair.filename)
		if got != expected {
			t.Errorf(
				"Basename('%s', '%s'): expected %s, got %s",
				pair.prefix,
				pair.filename,
				expected,
				got,
			)
		} else {
			t.Logf(
				"Basename('%s', '%s'): %s (OK)",
				pair.prefix,
				pair.filename,
				got,
			)
		}
	}
}

func TestFilenames(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("using tempDir %s", tempDir)
	defer os.RemoveAll(tempDir)

	expected := map[string]bool{
		"foo.txt":            false,
		"bar/baz.dat":        false,
		"a/b/c/d/e.document": false,
	}

	for file, _ := range expected {
		file = tempDir + "/" + file
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			t.Fatal(err)
		}
		f, err := os.Create(file)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	for _, suffix := range []string{"", "/"} {
		got := Filenames(tempDir + suffix)
		if len(got) != len(expected) {
			t.Fatalf("Filenames: expected %d, got %d", len(expected), len(got))
		}

		for _, gotFile := range got {
			if _, ok := expected[gotFile]; ok {
				expected[gotFile] = true
				t.Logf("'%s': expected and received", gotFile)
			} else {
				t.Errorf("'%s': got it, but didn't expect it", gotFile)
			}
		}

		for expectedFile, verified := range expected {
			if !verified {
				t.Errorf("'%s': expected it, but didn't get it", expectedFile)
			} else {
				t.Logf("'%s': verified receipt", expectedFile)
			}
		}
	}
}
