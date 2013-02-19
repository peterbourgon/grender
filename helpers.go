package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

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

func copyFile(dst, src string) error {
	return fmt.Errorf("not yet implemented")
}

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
