package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// splitPath tokenizes the given path string on filepath.Separator.
func splitPath(path string) []string {
	list := []string{}
	for {
		dir, file := filepath.Split(path)
		if file == "" {
			break
		}
		list = append([]string{file}, list...)
		if dir == "" {
			break
		}
		path = filepath.Clean(dir)
	}
	return list
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

// copyFile copies src to dst.
func copyFile(dst, src string) {
	mustWrite(dst, mustRead(src))
}

// mustRead returns the content of the passed filename.
func mustRead(filename string) []byte {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("%s: %s", filename, err)
		os.Exit(1)
	}
	return buf
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

// writeTo writes the buffer to the target file.
func mustWrite(tgt string, buf []byte) {
	os.MkdirAll(filepath.Dir(tgt), 0777)
	if err := ioutil.WriteFile(tgt, buf, 0755); err != nil {
		log.Printf("%s: %s", tgt, err)
		os.Exit(1)
	}
}

// targetFor returns the target filename for the given source filename.
func targetFor(filename string) string {
	dst := filepath.Clean(filepath.Join(*targetDir, diffPath(*sourceDir, filename)))
	n := len(dst) - len(filepath.Ext(dst))
	return dst[:n] + ".html"
}
