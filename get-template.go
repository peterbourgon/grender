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
	templateFile := ""
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
		if ValidTemplate(path, pseudoPageFile) {
			Debugf("  + VALID")
			templateFile = path
			return nil
		}
		Debugf("  · IGNORE")
		return nil
	}
	Debugf("walking %s", templateRoot)
	if err := filepath.Walk(templateRoot, w); err != nil {
		Problemf("Walk: %s", err)
	}
	if templateFile == "" {
		return "", fmt.Errorf("no matching template found")
	}
	Debugf("%s: chose %s", pageFile, templateFile)
	return templateFile, nil
}

func ValidTemplate(file, pageFile string) bool {
	switch path.Ext(file) {
	case TemplateExtension:
		break
	default:
		return false
	}
	if StripExtension(file) == StripExtension(pageFile) {
		return true // foo/bar/baz.html + foo/bar/baz.page => good
	}
	if path.Base(file) == "_"+TemplateExtension {
		return true // global environment (we only walk to the ones that apply)
	}
	return false
}

func Rehome(pageFile, sourceRoot, templateRoot string) string {
	return strings.Replace(pageFile, sourceRoot, templateRoot, 1)
}
