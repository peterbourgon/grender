package main

import (
	"testing"
	"io/ioutil"
	"os"
	"fmt"
)

func TestTokenizePath(t *testing.T) {
	m := map[string][]string{
		"/a/b/c":   []string{"a", "b", "c"},
		"./foo":    []string{"foo"},
		"x_y/z_a/": []string{"x_y", "z_a"},
	}
	for path, expected := range m {
		if got := TokenizePath(path); !equal(got, expected) {
			t.Errorf("%s: expected '%s', got '%s'", path, expected, got)
		}
	}
}

func TestSubpath(t *testing.T) {
	type pair struct {
		rootDir string
		file    string
	}
	m := map[pair]string{
		pair{"sub", "sub/dir"}:                         "dir",
		pair{"sub", "sub/dir/foo/"}:                    "dir/foo",
		pair{"sub", "sub/dir/foo/bar"}:                 "dir/foo/bar",
		pair{"sub/dir", "sub/dir/foo/bar"}:             "foo/bar",
		pair{"sub/dir/foo", "sub/dir/foo/bar"}:         "bar",
		pair{"sub/dir/foo/", "sub/dir/foo/bar"}:        "bar",
		pair{"/var/log", "/var/log/daemon.log"}:        "daemon.log",
		pair{"/usr///share/", "/usr/share/dict/words"}: "dict/words",
	}
	for pair, expected := range m {
		if got := Subpath(pair.rootDir, pair.file); got != expected {
			t.Errorf("%s: expected '%s', got '%s'", pair.file, expected, got)
		}
	}
}

func TestStripExtension(t *testing.T) {
	m := map[string]string{
		"foo.bar":             "foo",
		"A.long-extension":    "A",
		"beep":                "beep",
		"/some/path/file.txt": "/some/path/file",
		"./foo.x":             "./foo",
	}
	for file, expected := range m {
		if got := StripExtension(file); got != expected {
			t.Errorf("%s: expected '%s', got '%s'", file, expected, got)
		}
	}
}

func writeToTempFile(t *testing.T, contents string) string {
	f, err := ioutil.TempFile(os.TempDir(), "grender-context-test")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	if _, err := f.WriteString(contents); err != nil {
		t.Fatalf("%s", err)
	}
	return f.Name()
}

func TestMergeJSON(t *testing.T) {
	m := map[string]Context{
		`{}`: Context{},
		`{"a": "b"}`: Context{
			"a": "b",
		},
		`{"array": [1,2,3]}`: Context{
			"array": []int{1, 2, 3},
		},
		`{"map": {"a":1.23}}`: Context{
			"map": map[string]float64{
				"a": 1.23,
			},
		},
		`{"x": ["a", {"b":1}, "c"]}`: Context{
			"x": []interface{}{
				"a",
				map[string]int{
					"b": 1,
				},
				"c",
			},
		},
	}
	for body, expected := range m {
		filename := writeToTempFile(t, body)
		defer os.Remove(filename)
		got := Context{}
		mergeJSON(filename, got)
		if len(expected) != len(got) {
			t.Errorf("%s: expected %d, got %d", body, len(expected), len(got))
		}
		expectedStr := fmt.Sprintf("%v", expected)
		gotStr := fmt.Sprintf("%v", got)
		t.Logf("%s: expected '%s', got '%s'", body, expectedStr, gotStr)
		if expectedStr != gotStr {
			t.Errorf("%s: expected '%s', got '%s'", body, expectedStr, gotStr)
		}
	}
}

func TestMergeMarkdown(t *testing.T) {
	m := map[string]Context{
		``:  Context{MarkdownKey: "\n"},
		md1: Context{MarkdownKey: html1},
	}
	for body, expected := range m {
		filename := writeToTempFile(t, body)
		defer os.Remove(filename)
		got := Context{}
		mergeMarkdown(filename, got)
		if len(expected) != len(got) {
			t.Errorf("%s: expected %d, got %d", body, len(expected), len(got))
		}
		expectedStr := fmt.Sprintf("%v", expected)
		gotStr := fmt.Sprintf("%v", got)
		if expectedStr != gotStr {
			t.Errorf("%s: expected '%s', got '%s'", body, expectedStr, gotStr)
		}
	}
}

//
//
//

const md1 = `# A
Lorem _ipsum_ **dolor** sit.

Herpa derp—a derp.

## B

AT&T.
`

const html1 = `<h1>A</h1>

<p>Lorem <em>ipsum</em> <strong>dolor</strong> sit.</p>

<p>Herpa derp—a derp.</p>

<h2>B</h2>

<p>AT&amp;T.</p>
`
