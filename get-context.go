package main

import (
	"bytes"
	"os"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"path"
	"strings"
	"github.com/knieriem/markdown"
	. "github.com/peterbourgon/bonus/xlog"
)

const (
	JSONExtension     = ".json"
	MarkdownExtension = ".md"
	RawdataExtension  = ".rawdata"
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
			Infof("+ %s: merging data from %s", pageFile, path)
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
	Debugf("### ShouldMerge('%s', '%s')", file, pageFile)
	fullName, ext := SplitAtExtension(file)
	name := ""
	switch ext {
	case JSONExtension, MarkdownExtension:
		name = path.Base(fullName)
	case RawdataExtension:
		// these have the special naming sceheme: baz.keyname.rawdata
		// so to pass the below test, we need to strip the 1st extension
		name, _ = SplitAtExtension(path.Base(fullName))
	default:
		Debugf("### ShouldMerge('%s', '%s') NO 1", file, pageFile)
		return false
	}
	if fullName == StripExtension(pageFile) {
		Debugf("### ShouldMerge('%s', '%s') YES 1", file, pageFile)
		return true // foo/bar/baz.json + foo/bar/baz.page => merge
	}
	if name == "_" {
		Debugf("### ShouldMerge('%s', '%s') YES 2", file, pageFile)
		return true
	}
	Debugf("### ShouldMerge('%s', '%s') NO X fullName='%s'", file, pageFile, name)
	return false
}

func Merge(file string, ctx Context) {
	switch strings.ToLower(path.Ext(file)) {
	case JSONExtension:
		mergeJSON(file, ctx)
	case MarkdownExtension:
		mergeMarkdown(file, ctx)
	case RawdataExtension:
		mergeRawdata(file, ctx)
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

func mergeRawdata(file string, ctx Context) {
	f, err := os.Open(file)
	if err != nil {
		Problemf("%s: %s", file, err)
		return
	}
	defer f.Close()
	b1, e1 := SplitAtExtension(file)
	if e1 != RawdataExtension {
		Problemf("%s: %s != %s", file, e1, RawdataExtension)
		return
	}
	_, e2 := SplitAtExtension(b1)
	key := strings.TrimLeft(e2, ".")
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		Problemf("%s: %s", file, err)
		return
	}
	Debugf("%s: '%s' = %d bytes", file, key, len(buf))
	ctx[key] = string(buf)
	Debugf("after merging %s, ctx: %v", file, ctx)
}
