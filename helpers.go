package main

import (
	"encoding/json"
	"fmt"
	"github.com/peterbourgon/mergemap"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Read returns the content of the passed filename.
func Read(filename string) []byte {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		Fatalf("must read: %s: %s", filename, err)
	}
	return buf
}

// Write writes the buffer to the target file.
func Write(tgt string, buf []byte) {
	os.MkdirAll(filepath.Dir(tgt), 0777)
	if err := ioutil.WriteFile(tgt, buf, 0755); err != nil {
		Fatalf("must write: %s: %s", tgt, err)
	}
}

// Relative gives the relative path from base for complete. complete must have
// base as a prefix.
func Relative(base, complete string) string {
	rel, err := filepath.Rel(base, complete)
	if err != nil {
		Fatalf("Diff Path %s, %s: %s", base, complete, err)
	}

	// special case
	if rel == "." {
		rel = ""
	}

	return rel
}

// Copy copies src to dst.
func Copy(dst, src string) {
	Write(dst, Read(src))
}

// ParseJSON parses the passed JSON buffer and returns a map.
func ParseJSON(buf []byte) map[string]interface{} {
	m := map[string]interface{}{}
	if err := json.Unmarshal(buf, &m); err != nil {
		Fatalf("parse JSON: %s", err)
	}
	return m
}

// TargetFor returns the target filename for the given source filename.
func TargetFor(sourceFilename, targetExt string) string {
	relativePath := Relative(*sourceDir, sourceFilename)
	dst := filepath.Clean(filepath.Join(*targetDir, relativePath))
	n := len(dst) - len(filepath.Ext(dst))
	return dst[:n] + targetExt
}

// MaybeTemplate returns the contents of the template file specified under the
// "template" key for the metadata in the stack identified by the given path.
// In human words, it means "get me the template for this file".
func MaybeTemplate(s StackReader, path string) (string, []byte, error) {
	templateInterface, ok := s.Get(path)["template"]
	if !ok {
		return "", []byte{}, fmt.Errorf("%s: no template", path)
	}
	templateStr, ok := templateInterface.(string)
	if !ok {
		return "", []byte{}, fmt.Errorf("%s: bad type for template key", path)
	}
	templateFilename := filepath.Join(filepath.Dir(path), templateStr) // rel
	return templateFilename, Read(templateFilename), nil
}

// Template calls MaybeTemplate, and fatals on error.
func Template(s StackReader, path string) (string, []byte) {
	path, buf, err := MaybeTemplate(s, path)
	if err != nil {
		Fatalf("must template: %s", err)
	}
	return path, buf
}

// SplitPath tokenizes the given path string on filepath.Separator.
func SplitPath(path string) []string {
	path = filepath.Clean(path)
	if path == "." {
		return []string{} // special-case; see TestSplitPath
	}

	sep := string(filepath.Separator)
	list := []string{}
	for _, s := range strings.Split(path, sep) {
		if s := strings.TrimSpace(s); s != "" {
			list = append(list, s)
		}
	}
	return list
}

var (
	BlogEntryRegexp = regexp.MustCompile(`^([0-9]+)-([0-9]+)-([0-9]+)-(.*)\.[^\.]*$`)
)

type BlogTuple struct {
	Year     int
	Month    int
	Day      int
	Title    string
	Filename string
}

func NewBlogTuple(path, targetExt string) (BlogTuple, bool) {
	path = filepath.Base(path)
	m := BlogEntryRegexp.FindAllStringSubmatch(path, -1)

	if len(m) <= 0 {
		Debugf("Blog Tuple: %s: failed to parse stage 0", path)
		return BlogTuple{}, false
	}

	if len(m[0]) < 5 {
		Debugf("Default Date: %s: failed to parse stage 1", path)
		return BlogTuple{}, false
	}

	if len(m[0][1]) <= 0 || len(m[0][2]) <= 0 || len(m[0][3]) <= 0 {
		Debugf("Default Date: %s: failed to parse stage 2", path)
		return BlogTuple{}, false
	}

	yyyy, err := strconv.ParseInt(m[0][1], 10, 32)
	if err != nil {
		Debugf("Default Date: %s: bad year '%s'", path, m[1])
		return BlogTuple{}, false
	}

	mm, err := strconv.ParseInt(m[0][2], 10, 32)
	if err != nil {
		Debugf("Default Date: %s: bad month '%s'", path, m[2])
		return BlogTuple{}, false
	}

	dd, err := strconv.ParseInt(m[0][3], 10, 32)
	if err != nil {
		Debugf("Default Date: %s: bad day '%s'", path, m[3])
		return BlogTuple{}, false
	}

	if len(m[0][4]) <= 0 {
		Debugf("Default Title: %s: failed to parse stage 3", path)
		return BlogTuple{}, false
	}

	filename := m[0][4]
	title := filename
	title = strings.Replace(title, "-", " ", -1)
	title = strings.Replace(title, "_", " ", -1)
	title = strings.ToTitle(string(title[0])) + title[1:]

	return BlogTuple{
		Year:     int(yyyy),
		Month:    int(mm),
		Day:      int(dd),
		Title:    title,
		Filename: filename + targetExt,
	}, true
}

func (bt BlogTuple) DateString() string {
	return fmt.Sprintf("%04d %02d %02d", bt.Year, bt.Month, bt.Day)
}

func (bt BlogTuple) TargetFile(baseDir string) string {
	return filepath.Join(
		baseDir,
		fmt.Sprintf("%04d", bt.Year),
		fmt.Sprintf("%02d", bt.Month),
		fmt.Sprintf("%02d", bt.Day),
		bt.Filename,
	)
}

func (bt BlogTuple) Redirects(baseDir string) []string {
	uniques := map[string]struct{}{}
	for _, yearFmt := range []string{"%d", "%04d"} {
		for _, monthFmt := range []string{"%d", "%02d"} {
			for _, dayFmt := range []string{"%d", "%02d"} {
				uniques[filepath.Join(
					baseDir,
					fmt.Sprintf(yearFmt, bt.Year),
					fmt.Sprintf(monthFmt, bt.Month),
					fmt.Sprintf(dayFmt, bt.Day),
					fmt.Sprintf(bt.Filename),
				)] = struct{}{}
				uniques[filepath.Join(
					baseDir,
					fmt.Sprintf(yearFmt, bt.Year),
					fmt.Sprintf(monthFmt, bt.Month),
					fmt.Sprintf(dayFmt, bt.Day),
					"index.html",
				)] = struct{}{}
			}
		}
	}
	delete(uniques, bt.TargetFile(baseDir))

	redirects := []string{}
	for unique := range uniques {
		redirects = append(redirects, unique)
	}
	return redirects
}

func RedirectTo(url string) []byte {
	return []byte(fmt.Sprintf(
		`
		<html><head>
		<meta http-equiv="refresh" content="0;url=%s">
		</head><body></body></html>
		`,
		url,
	))
}

// SplatInto splits the `path` on filepath.Separator, and merges the passed
// `metadata` into the map `m` under the resulting key.
//
// As an example, if path="foo/bar/baz", SplatInto is semantically the same as
// `m = merge(m[foo][bar][baz], metadata)`.
func SplatInto(m map[string]interface{}, path string, metadata map[string]interface{}) {
	m0 := m
	for _, level := range SplitPath(path) {
		if _, ok := m0[level]; !ok {
			m0[level] = map[string]interface{}{}
		}
		m0 = m0[level].(map[string]interface{})
	}

	m0 = mergemap.Merge(m0, metadata)
}

func PrettyPrint(i interface{}) string {
	buf, _ := json.MarshalIndent(i, "# ", "    ")
	return string(buf)
}

// SortedValues returns a slice of every value in the passed map, ordered by
// the "sortkey" (if it exists) or the name of the entry (if it doesn't).
func SortedValues(m map[string]interface{}) []interface{} {
	mapping := map[string]string{} // sort key: original key
	for name, element := range m {
		submap, ok := element.(map[string]interface{})
		if !ok {
			mapping[name] = name
			continue
		}
		sortkey, ok := submap["sortkey"]
		if !ok {
			mapping[name] = name
			continue
		}
		sortkeyString, ok := sortkey.(string)
		if !ok {
			mapping[name] = name
			continue
		}
		mapping[sortkeyString] = name
	}

	sortkeys := stringSlice{}
	for sortkey, _ := range mapping {
		sortkeys = append(sortkeys, sortkey)
	}
	sort.Sort(sortkeys)

	orderedValues := []interface{}{}
	for _, k := range sortkeys {
		orderedValues = append(orderedValues, m[mapping[k]])
	}
	return orderedValues
}

type stringSlice []string

func (a stringSlice) Len() int           { return len(a) }
func (a stringSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a stringSlice) Less(i, j int) bool { return a[i] > a[j] }
