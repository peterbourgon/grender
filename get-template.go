package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	. "github.com/peterbourgon/bonus/xlog"
)

const (
	TemplateExtension = ".html"
)

func GetTemplate(sourceRoot, templateRoot, pageFile string) (string, error) {
	// pageFile is within sourceRoot
	pseudoPageFile := Rehome(pageFile, sourceRoot, templateRoot)
	Debugf("Rehome(%s, %s, %s) -> %s", pageFile, sourceRoot, templateRoot, pseudoPageFile)
	bestTemplate, bestValidity := "", 0
	w := func(path string, info os.FileInfo, err error) error {
		Debugf("• %s (%v)", path, err)
		if err != nil {
			Debugf("  X ERROR %s", err)
			return err
		}
		if info.IsDir() {
			Debugf("  · ShouldDescend(%s, %s) -> %s", path, pseudoPageFile, ShouldDescend(path, pseudoPageFile))
			if ShouldDescend(path, pseudoPageFile) {
				Debugf("  > DESCEND")
				return nil
			}
			Debugf("  - SKIPDIR")
			return filepath.SkipDir
		}
		Debugf("ValidTemplate(%s, %s) = %v", path, pseudoPageFile, ValidTemplate(path, pseudoPageFile))
		if validity := ValidTemplate(path, pseudoPageFile); validity > 0 {
			if validity > bestValidity {
				Debugf("  + VALID + CHOSEN")
				bestTemplate = path
				bestValidity = validity
				return nil
			}
			Debugf("  · VALID, NOT CHOSEN")
			return nil
		}
		Debugf("  · INVALID, IGNORE")
		return nil
	}
	Debugf("walking %s", templateRoot)
	if err := DepthFirstWalk(templateRoot, w); err != nil {
		Problemf("Walk: %s", err)
	}
	if bestTemplate == "" {
		return "", fmt.Errorf("no matching template found")
	}
	Debugf("%s: chose %s", pageFile, bestTemplate)
	return bestTemplate, nil
}

func ValidTemplate(file, pageFile string) int {
	switch path.Ext(file) {
	case TemplateExtension:
		break
	default:
		return 0
	}
	if StripExtension(file) == StripExtension(pageFile) {
		return 2 // foo/bar/baz.html + foo/bar/baz.page => good
	}
	if path.Base(file) == "_"+TemplateExtension {
		return 1 // global environment (we only walk to the ones that apply)
	}
	return 0
}

func Rehome(pageFile, sourceRoot, templateRoot string) string {
	return strings.Replace(pageFile, sourceRoot, templateRoot, 1)
}
