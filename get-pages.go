package main

import (
	"os"
	"path"
	"path/filepath"
	. "github.com/peterbourgon/bonus/xlog"
)

func GetPages(sourceRoot string) []string {
	pages := []string{}
	w := func(file string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return filepath.SkipDir
		}
		if path.Ext(file) == PageExtension {
			pages = append(pages, file)
		}
		return nil
	}
	if err := filepath.Walk(sourceRoot, w); err != nil {
		Problemf("%s: %s", sourceRoot, err)
		return []string{}
	}
	return pages
}
