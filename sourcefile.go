package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	MaxGetCount = 999999999
	YYYYMMDDT   = `([0-9]{4})-([0-9]{2})-([0-9]{2})-([0-9A-Za-z_\-\.]+)`
)

var (
	R = regexp.MustCompile(YYYYMMDDT)
)

// SourceFile is a representation of a parsed Source File,
// with the important bits explicitly extracted.
type SourceFile struct {
	SourceFile string
	Basename   string
	Metadata   map[string]interface{} // user-supplied metadata
}

type SourceFiles []*SourceFile

func NewSourceFile(filename string) *SourceFile {
	return &SourceFile{
		SourceFile: filename,
		Basename:   Basename(*sourcePath, filename),
		Metadata:   map[string]interface{}{},
	}
}

func parseBlogEntryRegex(basename string) ([]string, error) {
	a := R.FindAllStringSubmatch(basename, -1)
	if a == nil || len(a) <= 0 || len(a[0]) != 5 {
		return nil, fmt.Errorf("not a blog entry")
	}

	toks := a[0][1:]
	toks[3] = strings.Replace(toks[3], "-", " ", -1)
	if len(toks[3]) > 1 {
		toks[3] = strings.ToTitle(toks[3])[:1] + toks[3][1:]
	}

	return toks, nil
}

func (sf *SourceFile) Indexable() bool {
	if _, err := parseBlogEntryRegex(sf.Basename); err != nil {
		return false
	}
	return true
}

func (sf *SourceFile) BlogEntry() (y, m, d, t string, err error) {
	var a []string
	a, err = parseBlogEntryRegex(sf.Basename)
	if err != nil {
		return
	}

	y, m, d, t = a[0], a[1], a[2], a[3]
	return
}

func (sf *SourceFile) getAbstract(key string) interface{} {
	i, ok := sf.Metadata[key]
	if !ok {
		return nil
	}
	return i
}

func (sf *SourceFile) getCount(key string) (int, error) {
	i := sf.getAbstract(key)
	if b, ok := i.(bool); ok {
		if b {
			return MaxGetCount, nil
		} else {
			return 0, nil
		}
	}
	if n, ok := i.(int); ok {
		return n, nil
	}
	return 0, fmt.Errorf("failed to parse '%s' as count", key)
}

func (sf *SourceFile) getString(key string) string {
	s, ok := sf.getAbstract(key).(string)
	if !ok {
		return ""
	}
	return s
}

func (sf *SourceFile) Render() map[string]interface{} {
	return sf.Metadata
}
