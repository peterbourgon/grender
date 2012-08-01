package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	YYYYMMDD = "([0-9]{4})-([0-9]{2})-([0-9]{2})"
	Title    = "([0-9A-Za-z_-]+)"
)

var (
	R = regexp.MustCompile(fmt.Sprintf("%s-?%s?", YYYYMMDD, Title))
)

// SourceFile is a representation of a parsed Source File,
// with the important bits explicitly extracted.
type SourceFile struct {
	SourceFile   string
	Metadata     map[string]interface{} // user-supplied metadata
	TemplateFile string
	OutputFile   string
	IndexTuple   IndexTuple
	Content      string
}

func NewSourceFile(filename string) *SourceFile {
	return &SourceFile{
		SourceFile: filename,
		Metadata:   map[string]interface{}{},
	}
}

// IndexTuple is a representation of all fields in an index-tuple.
type IndexTuple struct {
	Type    string
	SortKey string
	Year    string
	Month   string
	Day     string
	Title   string
	URL     string
}

// ContributeTo merges the index-tuple to the global Index.
// It also Sorts the Index type it's being contributed to.
func (it *IndexTuple) ContributeTo(idx Index) {
	if targetArray, ok := idx[it.Type]; ok {
		idx[it.Type] = append(targetArray, it)
	} else {
		idx[it.Type] = OrderedIndexTuples{it}
	}
	sort.Sort(idx[it.Type])
}

// Render compiles the index-tuple down to a flat map[string]string.
func (it *IndexTuple) Render() map[string]string {
	m := map[string]string{
		*indexTupleTypeKey:    it.Type,
		*indexTupleSortKeyKey: it.SortKey,
		"year":                it.Year,
		"month":               it.Month,
		"day":                 it.Day,
		"url":                 it.URL,
	}
	if it.Title != "" {
		m["title"] = it.Title
	}
	return m
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
	basename := Basename(*sourcePath, filename)
	if y, m, d, t, err := blogEntry(basename); err == nil {
		sf.OutputFile = fmt.Sprintf("%s/%s", *blogPath, basename)

		sf.IndexTuple.Type = *defaultIndexTupleType
		sf.IndexTuple.SortKey = basename
		sf.IndexTuple.Year = y
		sf.IndexTuple.Month = m
		sf.IndexTuple.Day = d
		sf.IndexTuple.Title = t
		sf.IndexTuple.URL = fmt.Sprintf("%s.%s", sf.OutputFile, *outputExtension)
	}

	// read remaining metadata as YAML
	if err = goyaml.Unmarshal(buf, sf.Metadata); err != nil {
		return
	}

	// index-tuple related metadata gets copied over
	if m0, ok := sf.Metadata[*indexTupleKey]; ok {
		if m1, ok := m0.(map[interface{}]interface{}); ok {
			for k, v := range m1 {
				kStr, kOk := k.(string)
				if !kOk {
					continue
				}
				vStr, vOk := v.(string)
				if !vOk {
					continue
				}
				copyMetadata(&sf.IndexTuple, strings.ToLower(kStr), vStr)
			}
		}
	}

	// check for template key: missing = fatal
	i, ok := sf.Metadata[*templateKey]
	if !ok {
		err = fmt.Errorf("%s: '%s' not provided", filename, *templateKey)
		return
	}
	if sf.TemplateFile, ok = i.(string); !ok {
		err = fmt.Errorf("%s: '%s' not a string", filename, *templateKey)
		return
	}

	// check for output file key: missing = need to deduce from basename
	if sf.OutputFile == "" {
		i, ok = sf.Metadata[*outputKey]
		if ok {
			if sf.OutputFile, ok = i.(string); !ok {
				err = fmt.Errorf("%s: '%s' not a string", filename, *outputKey)
				return
			}
		} else {
			sf.OutputFile = Basename(*sourcePath, filename)
		}
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

func copyMetadata(it *IndexTuple, k, v string) {
	switch k {
	case *indexTupleTypeKey:
		it.Type = v
	case *indexTupleSortKeyKey:
		it.SortKey = v
	case "year":
		it.Year = v
	case "month":
		it.Month = v
	case "day":
		it.Day = v
	case "title":
		it.Title = v
	case "url":
		it.URL = v
	}
}
