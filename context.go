package main

import (
	"os"
	"path/filepath"
	"path"
	"strings"
	. "github.com/peterbourgon/bonus/xlog"
)

type Node struct {
	file	[]string
	subdirs	[]Node
}

type Context map[string]interface{}

func GetContext(sourceRoot, templateRoot, pageName string) Context {
	ctx := Context{}
	w := func(path string, info os.FileInfo, err error) error {
		Debugf("• %s", path)
		if info.IsDir() {

		}
		return nil
	}
	Debugf("walking %s", sourceRoot)
	if err := filepath.Walk(sourceRoot, w); err != nil {
		Problemf("Walk: %s", err)
	}
	return ctx
}

func ShouldDescend(dir, pageFile string) bool {
	d, pf := TokenizePath(dir), TokenizePath(pageFile)
	if len(d) >= len(pf) {
		return false
	}
	if equal(d[:len(pf)], pf) {
		return true
	}
	return false
}

func ShouldMerge(file, pageFile string) bool {
	if StripExtension(file) == StripExtension(pageFile) {
		return true	// foo/bar/baz.json + foo/bar/baz.page => merge
	}
	if path.Base(file) == "_.json" {
		return true	// global environment (we only walk to the ones that apply)
	}
	return false
}

func Subpath(rootDir, file string) string {
	d, err := filepath.Abs(rootDir)
	if err != nil {
		Problemf("PathFor(%s, %s): %s", rootDir, file, err)
		return ""
	}
	f, err := filepath.Abs(file)
	if err != nil {
		Problemf("PathFor(%s, %s): %s", rootDir, file, err)
		return ""
	}
	if strings.Index(f, d) != 0 {
		Problemf("PathFor(%s, %s): file not under directory", rootDir, file)
		return ""
	}
	return f[len(d)+1:]
}

func TokenizePath(s string) []string {
	p, a := s, []string{}
	for {
		p = path.Clean(p)
		l, r := path.Dir(p), path.Base(p)
		if r == "." || r == "/" {
			break
		}
		a = append(a, r)
		p = l
	}
	reverse(a)
	return a
}

func StripExtension(file string) string {
	return file[:len(file)-len(path.Ext(file))]
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func reverse(a []string) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
