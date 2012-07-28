package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBadFile(t *testing.T) {
	badFilename := "/no/such/file"
	if _, err := ParseSourceFile(
		"",
		badFilename,
		*metadataDelimiter,
		*contentKey,
		*templateKey,
		*outputKey,
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
		*metadataDelimiter,
		*contentKey,
		*templateKey,
		*outputKey,
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
			*metadataDelimiter,
			*contentKey,
			*templateKey,
			*outputKey,
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
