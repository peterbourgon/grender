package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
// base as a prefix; otherwise, diffPath returns complete, unaltered.
func diffPath(base, complete string) string {
	base, complete = filepath.Clean(base), filepath.Clean(complete)

	if len(complete) <= len(base) {
		return complete
	}
	if complete[:len(base)] != base {
		return complete
	}

	d := complete[len(base):]
	if d[0] == filepath.Separator {
		d = d[1:]
	}
	return d
}

// copyFile copies src to dst using primitives from the os package.
func copyFile(dst, src string) error {
	return fmt.Errorf("not yet implemented")
}

// readJSON interprets the content of filename as JSON data, and returns a
// map of that data or an error.
func readJSON(filename string) (map[string]interface{}, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(buf, &m); err != nil {
		return nil, err
	}

	return m, nil
}
