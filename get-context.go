package main

import (
	"bytes"
	"os"
	"encoding/json"
	"path/filepath"
	"path"
	"strings"
	"github.com/knieriem/markdown"
	. "github.com/peterbourgon/bonus/xlog"
)

const (
	JSONExtension     = ".json"
	MarkdownExtension = ".md"
	PageExtension     = ".page"
	MarkdownKey       = "markdown"
)

type Context map[string]interface{}

func GetContext(sourceRoot, pageFile string) Context {
	ctx := Context{}
	w := func(path string, info os.FileInfo, err error) error {
		Debugf("• %s (%v)", path, err)
		if err != nil {
			Debugf("  X ERROR %s", err)
			return err
		}
		if info.IsDir() {
			if ShouldDescend(path, pageFile) {
				Debugf("  > DESCEND")
				return nil
			}
			Debugf("  - SKIPDIR")
			return filepath.SkipDir
		}
		if ShouldMerge(path, pageFile) {
			Debugf("  + MERGE")
			Merge(path, ctx)
			return nil
		}
		Debugf("  · IGNORE")
		return nil
	}
	Debugf("walking %s", sourceRoot)
	if err := filepath.Walk(sourceRoot, w); err != nil {
		Problemf("Walk: %s", err)
	}
	return ctx
}

func ShouldMerge(file, pageFile string) bool {
	switch path.Ext(file) {
	case JSONExtension, MarkdownExtension:
		break
	default:
		return false
	}
	if StripExtension(file) == StripExtension(pageFile) {
		return true // foo/bar/baz.json + foo/bar/baz.page => merge
	}
	if path.Base(file) == "_"+JSONExtension {
		return true // global environment (we only walk to the ones that apply)
	}
	return false
}

func Merge(file string, ctx Context) {
	switch strings.ToLower(path.Ext(file)) {
	case JSONExtension:
		mergeJSON(file, ctx)
	case MarkdownExtension:
		mergeMarkdown(file, ctx)
	default:
		Problemf("%s: unknown type; not merging", file)
	}
}

func mergeJSON(file string, ctx Context) {
	f, err := os.Open(file)
	if err != nil {
		Problemf("%s: %s", file, err)
		return
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&ctx); err != nil {
		Problemf("%s: decode: %s", file, err)
	}
}

func mergeMarkdown(file string, ctx Context) {
	f, err := os.Open(file)
	if err != nil {
		Problemf("%s: %s", file, err)
		return
	}
	defer f.Close()
	p := markdown.NewParser(&markdown.Extensions{Smart: true})
	b := bytes.Buffer{}
	p.Markdown(f, markdown.ToHTML(&b))
	ctx[MarkdownKey] = b.String()
}
