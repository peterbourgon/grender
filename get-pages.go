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
		Debugf("• %s", file)
		if info.IsDir() {
			Debugf("  DIVE")
			return nil
		}
		if path.Ext(file) == PageExtension {
			Debugf("  VALID")
			pages = append(pages, file)
			return nil
		}
		Debugf("  IGNORED")
		return nil
	}
	Debugf("GetPages(%s) looking for Pages", sourceRoot)
	if err := filepath.Walk(sourceRoot, w); err != nil {
		Problemf("%s: %s", sourceRoot, err)
		return []string{}
	}
	return pages
}
