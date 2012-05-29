package main

import (
	"testing"
)

func TestTokenizePath(t *testing.T) {
	m := map[string][]string{
		"/a/b/c":	[]string{"a", "b", "c"},
		"./foo":	[]string{"foo"},
		"x_y/z_a/":	[]string{"x_y", "z_a"},
	}
	for path, expected := range m {
		if got := TokenizePath(path); !equal(got, expected) {
			t.Errorf("%s: expected '%s', got '%s'", path, expected, got)
		}
	}
}

func TestSubpath(t *testing.T) {
	type pair struct {
		rootDir	string
		file	string
	}
	m := map[pair]string{
		pair{"sub", "sub/dir"}:				"dir",
		pair{"sub", "sub/dir/foo/"}:			"dir/foo",
		pair{"sub", "sub/dir/foo/bar"}:			"dir/foo/bar",
		pair{"sub/dir", "sub/dir/foo/bar"}:		"foo/bar",
		pair{"sub/dir/foo", "sub/dir/foo/bar"}:		"bar",
		pair{"sub/dir/foo/", "sub/dir/foo/bar"}:	"bar",
		pair{"/var/log", "/var/log/daemon.log"}:	"daemon.log",
		pair{"/usr///share/", "/usr/share/dict/words"}:	"dict/words",
	}
	for pair, expected := range m {
		if got := Subpath(pair.rootDir, pair.file); got != expected {
			t.Errorf("%s: expected '%s', got '%s'", pair.file, expected, got)
		}
	}
}

func TestStripExtension(t *testing.T) {
	m := map[string]string{
		"foo.bar":		"foo",
		"A.long-extension":	"A",
		"beep":			"beep",
		"/some/path/file.txt":	"/some/path/file",
		"./foo.x":		"./foo",
	}
	for file, expected := range m {
		if got := StripExtension(file); got != expected {
			t.Errorf("%s: expected '%s', got '%s'", file, expected, got)
		}
	}
}
