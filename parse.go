package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	YYYYMMDD = "([0-9]{4})-([0-9]{2})-([0-9]{2})"
	T        = "([0-9A-Za-z_-]+)"
)

var (
	R = regexp.MustCompile(fmt.Sprintf("%s-%s", YYYYMMDD, T))
)

// SourceFile is a representation of a parsed Source File,
// with the important bits explicitly extracted.
type SourceFile struct {
	SourceFile string
	Basename   string
	Metadata   map[string]interface{} // user-supplied metadata
	Content    string
}

type SourceFiles []*SourceFile

func NewSourceFile(filename string) *SourceFile {
	return &SourceFile{
		SourceFile: filename,
		Basename:   Basename(*sourcePath, filename),
		Metadata:   map[string]interface{}{},
	}
}

func (sf *SourceFile) Indexable() bool {
	a := R.FindAllStringSubmatch(sf.Basename, -1)
	if a == nil || len(a) <= 0 || len(a[0]) <= 4 {
		return false
	}
	return true
}

func (sf *SourceFile) SortKey() string {
	return sf.getString("sortkey")
}

func (sf *SourceFile) Template() string {
	return sf.getString(*templateKey)
}

func (sf *SourceFile) Output() string {
	return sf.getString(*outputKey)
}

func (sf *SourceFile) getAbstract(key string) interface{} {
	i, ok := sf.Metadata[key]
	if !ok {
		return nil
	}
	return i
}

func (sf *SourceFile) getBool(key string) bool {
	b, ok := sf.getAbstract(key).(bool)
	if !ok {
		return false
	}
	return b
}

func (sf *SourceFile) getString(key string) string {
	s, ok := sf.getAbstract(key).(string)
	if !ok {
		return ""
	}
	return s
}

// ParseSourceFile reads the given filename (assumed to be a relative file under
// *sourcePath) and produces a parsed SourceFile object from its contents.
func ParseSourceFile(filename string) (sf *SourceFile, err error) {
	sf = NewSourceFile(filename)

	// read file
	f, err := os.Open(*sourcePath + "/" + filename)
	if err != nil {
		return
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	s := string(buf)

	// separate content
	if idx := strings.Index(s, *metadataDelimiter); idx >= 0 {
		delimiterCutoff := idx + len(*metadataDelimiter) + 1 // plus '\n'
		contentBuf := buf[delimiterCutoff:]

		switch strings.ToLower(filepath.Ext(filename)) {
		case ".md":
			contentBuf = RenderMarkdown(contentBuf)
		}

		sf.Content = strings.TrimSpace(string(contentBuf))
		buf = buf[:idx] // buf shall contain only metadata
	}

	// if the filename looks like a blog entry, autopopulate some metadata
	if y, m, d, t, err := blogEntry(sf.Basename); err == nil {
		sf.Metadata["output"] = fmt.Sprintf("%s/%s", *blogPath, sf.Basename)
		sf.Metadata["sortkey"] = sf.Basename
		sf.Metadata["year"] = y
		sf.Metadata["month"] = m
		sf.Metadata["day"] = d
		sf.Metadata["title"] = t
		sf.Metadata["url"] = fmt.Sprintf("%s.%s", sf.Output(), *outputExtension)
	}

	// read remaining metadata as YAML
	if err = goyaml.Unmarshal(buf, sf.Metadata); err != nil {
		return
	}

	// check for template key: missing = fatal
	if sf.Template() == "" {
		err = fmt.Errorf("%s: '%s' not provided", filename, *templateKey)
		return
	}

	// check for output file key: missing = need to deduce from basename
	if sf.Output() == "" {
		sf.Metadata[*outputKey] = Basename(*sourcePath, filename)
	}

	err = nil // just in case
	return
}

func blogEntry(basename string) (y, m, d, t string, err error) {
	a := R.FindAllStringSubmatch(basename, -1)
	if a == nil || len(a) <= 0 || len(a[0]) <= 3 {
		err = fmt.Errorf("not a blog entry")
		return
	}

	y, m, d, t = a[0][1], a[0][2], a[0][3], ""
	if len(a[0]) > 3 {
		t = strings.Replace(a[0][4], "-", " ", -1)
		if len(t) > 1 {
			t = strings.ToTitle(t)[:1] + t[1:]
		}
	}

	return
}
