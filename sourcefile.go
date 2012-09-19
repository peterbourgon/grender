package main

import (
	"fmt"
	"regexp"
	"strconv"
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
	toks[3] = Filename2Display(toks[3])
	return toks, nil
}

func Filename2Display(title string) string {
	title = strings.Replace(title, "-", " ", -1)
	if len(title) > 1 {
		title = strings.ToTitle(title)[:1] + title[1:]
	}
	return title
}

func Display2Filename(title string) string {
	title = strings.Replace(title, " ", "-", -1)
	title = strings.ToLower(title)
	return title
}

func (sf *SourceFile) Indexable() bool {
	if _, err := parseBlogEntryRegex(sf.Basename); err != nil {
		return false
	}
	return true
}

func (sf *SourceFile) BlogEntry() (y, m, d int, t string, err error) {
	a, err := parseBlogEntryRegex(sf.Basename)
	if err != nil {
		return 0, 0, 0, "", err
	}

	yStr, mStr, dStr := a[0], a[1], a[2]
	y64, err := strconv.ParseInt(yStr, 10, 32)
	if err != nil {
		return 0, 0, 0, "", err
	}
	m64, err := strconv.ParseInt(mStr, 10, 32)
	if err != nil {
		return 0, 0, 0, "", err
	}
	d64, err := strconv.ParseInt(dStr, 10, 32)
	if err != nil {
		return
	}

	return int(y64), int(m64), int(d64), a[3], nil
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
