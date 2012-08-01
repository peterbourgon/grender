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
	if _, err := ParseSourceFile(
		"",
		badFilename,
		NewIndex(),
		*blogPath,
		*metadataDelimiter,
		*contentKey,
		*templateKey,
		*indexKey,
		*outputKey,
		*outputExtension,
	); err == nil {
		t.Fatalf("ParseSourceFile successfully read %s", badFilename)
	} else {
		t.Logf("ParseSourceFile('%s') gave error: %s (good!)", badFilename, err)
	}
}

const simplestBody = `
template: nosuch.template
---
`

func writeTempFile(t *testing.T, filename, body string) (tempDir string) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("using tempDir %s", tempDir)

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

func TestAllKeys(t *testing.T) {
	tempFile := "dummy.src"
	tempDir := writeTempFile(t, tempFile, simplestBody)
	defer os.RemoveAll(tempDir)

	ctx, err := ParseSourceFile(
		tempDir,
		tempFile,
		NewIndex(),
		*blogPath,
		*metadataDelimiter,
		*contentKey,
		*templateKey,
		*indexKey,
		*outputKey,
		*outputExtension,
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{*contentKey, *templateKey, *outputKey} {
		if _, ok := ctx[key]; !ok {
			t.Errorf("missing '%s' key", key)
		} else {
			t.Logf("got '%s' key OK", key)
		}
	}
}

func TestDeducedOutputFilename(t *testing.T) {
	for _, sourceFilename := range []string{"foo.src", "a/b/c.txt"} {
		tempDir := writeTempFile(t, sourceFilename, simplestBody)
		defer os.RemoveAll(tempDir)

		ctx, err := ParseSourceFile(
			tempDir,
			sourceFilename,
			NewIndex(),
			*blogPath,
			*metadataDelimiter,
			*contentKey,
			*templateKey,
			*indexKey,
			*outputKey,
			*outputExtension,
		)
		if err != nil {
			t.Errorf("%s: parsing: %s", sourceFilename, err)
			continue
		}

		got, ok := ctx[*outputKey]
		if !ok {
			t.Errorf("%s: context missing '%s'", sourceFilename, *outputKey)
			continue
		}

		expected := Basename(tempDir, sourceFilename)
		if got != expected {
			t.Errorf("expected '%s', got '%s'", expected, got)
			continue
		}

		t.Logf("'%s' gave %s=%s OK", sourceFilename, *outputKey, got)
	}
}

func TestAutopopulatedIndex(t *testing.T) {
	m := map[string]string{
		"2012-01-01.md":             "",
		"2012-01-01-hello.md":       "Hello",
		"2012-01-01-hello-there.md": "Hello there",
	}

	for tempFile, expectedTitle := range m {
		tempDir := writeTempFile(t, tempFile, simplestBody)
		defer os.RemoveAll(tempDir)

		ctx, err := ParseSourceFile(
			tempDir,
			tempFile,
			NewIndex(),
			*blogPath,
			*metadataDelimiter,
			*contentKey,
			*templateKey,
			*indexKey,
			*outputKey,
			*outputExtension,
		)
		if err != nil {
			t.Fatal(err)
		}

		idxAbstract, ok := ctx[*indexKey]
		if !ok {
			t.Errorf("%s: no %s", tempFile, *indexKey)
			continue
		}

		idx, ok := idxAbstract.(map[string]string)
		if !ok {
			t.Errorf("%s: not a map-type", tempFile)
			continue
		}

		gotTitle, ok := idx["title"]
		if !ok {
			t.Errorf("%s: no %s", tempFile, "title")
			continue
		}

		if gotTitle != expectedTitle {
			t.Errorf("%s: got '%s', expected '%s'", tempFile, gotTitle, expectedTitle)
			continue
		}

		expectedURL := *blogPath + "/" + Basename("", tempFile) + "." + *outputExtension
		if gotURL := idx["url"]; gotURL != expectedURL {
			t.Errorf("%s: got '%s', expected '%s'", tempFile, gotURL, expectedURL)
			continue
		}
	}
}
