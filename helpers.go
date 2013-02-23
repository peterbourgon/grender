package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// mustRead returns the content of the passed filename.
func mustRead(filename string) []byte {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("%s: %s", filename, err)
		os.Exit(1)
	}
	return buf
}

// mustWrite writes the buffer to the target file.
func mustWrite(tgt string, buf []byte) {
	os.MkdirAll(filepath.Dir(tgt), 0777)
	if err := ioutil.WriteFile(tgt, buf, 0755); err != nil {
		log.Printf("%s: %s", tgt, err)
		os.Exit(1)
	}
}

// diffPath gives the relative path from base for complete. complete must have
// base as a prefix.
func diffPath(base, complete string) string {
	base, complete = filepath.Clean(base), filepath.Clean(complete)

	if len(complete) <= len(base) {
		log.Printf("diffPath('%s', '%s') invalid (length)", base, complete)
		os.Exit(1)
	}
	if complete[:len(base)] != base {
		log.Printf("diffPath('%s', '%s') invalid (prefix)", base, complete)
		os.Exit(1)
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
		log.Printf("%s", err)
		os.Exit(1)
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

// mustTemplate returns the contents of the template file specified under the
// 'template' key for the metadata in the stack identified by the given path.
// In human words, it means "get me the template for this file".
func mustTemplate(s *Stack, path string) []byte {
	template, ok := s.Get(path)["template"]
	if !ok {
		log.Printf("%s: no template", path)
		os.Exit(1)
	}
	templateStr, ok := template.(string)
	if !ok {
		log.Printf("%s: bad type for template key", path)
		os.Exit(1)
	}
	templateFile := filepath.Join(filepath.Dir(path), templateStr)
	return mustRead(templateFile)
}
