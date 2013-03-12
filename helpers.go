package main

import (
	"encoding/json"
	"fmt"
	"github.com/peterbourgon/mergemap"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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
		Fatalf("must JSON: %s", err)
	}
	return m
}

// TargetFor returns the target filename for the given source filename.
func TargetFor(sourceFilename, ext string) string {
	relativePath := Relative(*sourceDir, sourceFilename)
	dst := filepath.Clean(filepath.Join(*targetDir, relativePath))
	n := len(dst) - len(filepath.Ext(dst))
	return dst[:n] + ext
}

// MaybeTemplate returns the contents of the template file specified under the
// "template" key for the metadata in the stack identified by the given path.
// In human words, it means "get me the template for this file".
func MaybeTemplate(s StackReader, path string) ([]byte, error) {
	template, ok := s.Get(path)["template"]
	if !ok {
		return []byte{}, fmt.Errorf("%s: no template", path)
	}
	templateStr, ok := template.(string)
	if !ok {
		return []byte{}, fmt.Errorf("%s: bad type for template key", path)
	}
	templateFile := filepath.Join(*sourceDir, templateStr)
	return Read(templateFile), nil
}

// Template calls MaybeTemplate, and fatals on error.
func Template(s StackReader, path string) []byte {
	buf, err := MaybeTemplate(s, path)
	if err != nil {
		Fatalf("must template: %s", err)
	}
	return buf
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
