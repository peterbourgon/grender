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

// mustRead returns the content of the passed filename.
func mustRead(filename string) []byte {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		Fatalf("must read: %s: %s", filename, err)
	}
	return buf
}

// mustWrite writes the buffer to the target file.
func mustWrite(tgt string, buf []byte) {
	os.MkdirAll(filepath.Dir(tgt), 0777)
	if err := ioutil.WriteFile(tgt, buf, 0755); err != nil {
		Fatalf("must write: %s: %s", tgt, err)
	}
}

// diffPath gives the relative path from base for complete. complete must have
// base as a prefix.
func diffPath(base, complete string) string {
	base, complete = filepath.Clean(base), filepath.Clean(complete)

	if len(complete) < len(base) {
		Fatalf("diffPath('%s', '%s') invalid (length)", base, complete)
	}
	if complete[:len(base)] != base {
		Fatalf("diffPath('%s', '%s') invalid (prefix)", base, complete)
	}

	if len(complete) == len(base) {
		return ""
	}

	d := complete[len(base):]
	if d[0] == filepath.Separator {
		d = d[1:]
	}
	return d
}

// mustCopy copies src to dst.
func mustCopy(dst, src string) {
	mustWrite(dst, mustRead(src))
}

// readJSON parses the passed JSON buffer and returns a map.
func mustJSON(buf []byte) map[string]interface{} {
	m := map[string]interface{}{}
	if err := json.Unmarshal(buf, &m); err != nil {
		Fatalf("must JSON: %s", err)
	}
	return m
}

// targetFor returns the target filename for the given source filename.
func targetFor(sourceFilename, ext string) string {
	relativePath := diffPath(*sourceDir, sourceFilename)
	dst := filepath.Clean(filepath.Join(*targetDir, relativePath))
	n := len(dst) - len(filepath.Ext(dst))
	return dst[:n] + ext
}

// maybeTemplate returns the contents of the template file specified under the
// "template" key for the metadata in the stack identified by the given path.
// In human words, it means "get me the template for this file".
func maybeTemplate(s StackReader, path string) ([]byte, error) {
	template, ok := s.Get(path)["template"]
	if !ok {
		return []byte{}, fmt.Errorf("%s: no template", path)
	}
	templateStr, ok := template.(string)
	if !ok {
		return []byte{}, fmt.Errorf("%s: bad type for template key", path)
	}
	templateFile := filepath.Join(*sourceDir, templateStr)
	return mustRead(templateFile), nil
}

// mustTemplate calls maybeTemplate, and fatals on error.
func mustTemplate(s StackReader, path string) []byte {
	buf, err := maybeTemplate(s, path)
	if err != nil {
		Fatalf("must template: %s", err)
	}
	return buf
}

// splitPath tokenizes the given path string on filepath.Separator.
func splitPath(path string) []string {
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

// splatInto splits the `path` on filepath.Separator, and merges the passed
// `metadata` into the map `m` under the resulting key.
//
// As an example, if path="foo/bar/baz", splatInto is semantically the same as
// `m = merge(m[foo][bar][baz], metadata)`.
func splatInto(m map[string]interface{}, path string, metadata map[string]interface{}) {
	m0 := m
	for _, level := range splitPath(path) {
		if _, ok := m0[level]; !ok {
			m0[level] = map[string]interface{}{}
		}
		m0 = m0[level].(map[string]interface{})
	}

	m0 = mergemap.Merge(m0, metadata)
}

func prettyPrint(i interface{}) string {
	buf, _ := json.MarshalIndent(i, "# ", "    ")
	return string(buf)
}

// sortedValues returns a slice of every value in the passed map, ordered by
// the "sortkey" (if it exists) or the name of the entry (if it doesn't).
func sortedValues(m map[string]interface{}) []interface{} {
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
