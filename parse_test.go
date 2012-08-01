package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBlogEntryFilenameRegex(t *testing.T) {
	m := map[string]bool{
		"2012-01-01":              true,
		"2012-01-01-simple":       true,
		"2012-01-01-more-complex": true,
		"0000-00-00":              true,
		"1234-56-78-a-b-c_def":    true,
		"1234-56-78-_":            true,
		"1234-56-78-a_b_c":        true,
		"2012-1-1":                false,
		"2012-01-1":               false,
		"2012-1-01":               false,
	}
	for s, expected := range m {
		if got := R.MatchString(s); got != expected {
			t.Errorf("'%s': expected %s, got %s", s, expected, got)
		}
	}
}

func TestBadFile(t *testing.T) {
	badFilename := "/no/such/file"
	if _, err := ParseSourceFile(badFilename); err == nil {
		t.Fatalf("ParseSourceFile successfully read %s", badFilename)
	} else {
		t.Logf("ParseSourceFile('%s') gave error: %s (good!)", badFilename, err)
	}
}

const simplestBody = `
template: nosuch.template
---
`

func writeSourceFile(t *testing.T, filename, body string) (tempDir string) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("using tempDir %s", tempDir)
	*sourcePath = tempDir // impt. when checking eg. URL

	absTempFile := tempDir + "/" + filename
	if err := os.MkdirAll(filepath.Dir(absTempFile), 0755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(absTempFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := f.WriteString(body); err != nil {
		t.Fatal(err)
	}
	f.Close()

	return
}

func TestRequiredKeys(t *testing.T) {
	tempFile := "dummy.src"
	tempDir := writeSourceFile(t, tempFile, simplestBody)
	defer os.RemoveAll(tempDir)

	sf, err := ParseSourceFile(tempFile)
	if err != nil {
		t.Fatal(err)
	}

	if sf.TemplateFile == "" {
		t.Errorf("%s missing", *templateKey)
	}
	if sf.OutputFile == "" {
		t.Errorf("%s missing", *outputKey)
	}
}

func TestDeducedOutputFilename(t *testing.T) {
	for _, sourceFilename := range []string{"foo.src", "a/b/c.txt"} {
		tempDir := writeSourceFile(t, sourceFilename, simplestBody)
		defer os.RemoveAll(tempDir)

		sf, err := ParseSourceFile(sourceFilename)
		if err != nil {
			t.Errorf("%s: parsing: %s", sourceFilename, err)
			continue
		}

		got := sf.OutputFile
		expected := Basename(tempDir, sourceFilename)
		if got != expected {
			t.Errorf("expected '%s', got '%s'", expected, got)
			continue
		}
	}
}

func TestAutopopulatedIndexTupleTitles(t *testing.T) {
	m := map[string]string{
		"2012-01-01.md":             "",
		"2012-01-01-hello.md":       "Hello",
		"2012-01-01-hello-there.md": "Hello there",
	}

	for tempFile, expectedTitle := range m {
		tempDir := writeSourceFile(t, tempFile, simplestBody)
		defer os.RemoveAll(tempDir)

		sf, err := ParseSourceFile(tempFile)
		if err != nil {
			t.Fatal(err)
		}

		if gotTitle := sf.IndexTuple.Title; gotTitle != expectedTitle {
			t.Errorf("%s: got '%s', expected '%s'", tempFile, gotTitle, expectedTitle)
			continue
		}

		expectedURL := *blogPath + "/" + Basename("", tempFile) + "." + *outputExtension
		if gotURL := sf.IndexTuple.URL; gotURL != expectedURL {
			t.Errorf("%s: got '%s', expected '%s'", tempFile, gotURL, expectedURL)
			continue
		}
	}
}

const titledContent = `
template: nosuch.template
index:
   title: The INDEX TITLE!! from the Meta Data
---
Content of the thing.
`

func TestProperMergeOfIndexTupleMetadata(t *testing.T) {
	filename := "2012-01-01-test-proper-merge-of-index.md"
	tempDir := writeSourceFile(t, filename, titledContent)
	defer os.RemoveAll(tempDir)

	sf, err := ParseSourceFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	expectedTitle := "The INDEX TITLE!! from the Meta Data"
	if gotTitle := sf.IndexTuple.Title; gotTitle != expectedTitle {
		t.Fatalf("%s: got '%s', expected '%s'", filename, gotTitle, expectedTitle)
	}
}

func TestGlobalIndex(t *testing.T) {
	filename := "2012-01-01-testing-global-index.md"
	tempDir := writeSourceFile(t, filename, titledContent)
	defer os.RemoveAll(tempDir)

	sf, err := ParseSourceFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	idx := Index{}
	sf.IndexTuple.ContributeTo(idx)

	a, ok := idx[*defaultIndexTupleType]
	if !ok {
		t.Fatalf("%s not merged properly: '%s' missing", *indexTupleKey, *defaultIndexTupleType)
	}
	if len(a) != 1 {
		t.Fatalf("%s not merged properly: '%s' is len=%d", *indexTupleKey, *defaultIndexTupleType, len(a))
	}
	if a[0].Title != "The INDEX TITLE!! from the Meta Data" {
		t.Fatalf("%s not merged properly: bad title '%s'", *indexTupleKey, a[0].Title)
	}
}
