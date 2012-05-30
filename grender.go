package main

import (
	"flag"
	"github.com/peterbourgon/bonus/xlog"
	"github.com/hoisie/mustache"
	"strings"
	"os"
	"path"
	"fmt"
)

var (
	debug        *bool   = flag.Bool("debug", false, "enable debug logging")
	templatesDir *string = flag.String("templates-dir", "_templates", "directory containing templates")
	sourceDir    *string = flag.String("source-dir", "_source", "directory containing source")
	outputDir    *string = flag.String("output-dir", "_output", "directory where site will be written")
	temp         *string = flag.String("temp", "", "temporary flag")
)

const (
	OutputExtension = ".html"
)

func init() {
	flag.Parse()
	xlog.Initialize(*debug)
}

func main() {
	for _, page := range GetPages(*sourceDir) {
		xlog.Debugf("parsing %s", page)
		ctx := GetContext(*sourceDir, page)
		templateFile, err := GetTemplate(*sourceDir, *templatesDir, page)
		if err != nil {
			xlog.Problemf("%s: %s", page, err)
			continue
		}
		xlog.Infof("%s: chose template %s", page, templateFile)
		tmpl := mustache.RenderFile(templateFile, ctx)
		outputFile, err := WriteOutput(*sourceDir, *outputDir, page, tmpl)
		if err != nil {
			xlog.Problemf("%s: %s", page, err)
			continue
		}
		xlog.Infof("%s: wrote %d bytes to %s", page, len(tmpl), outputFile)
	}
}

func WriteOutput(sourceDir, outputDir, page, contents string) (string, error) {
	outputFile := strings.Replace(page, sourceDir, outputDir, 1)
	outputFile = strings.Replace(outputFile, PageExtension, OutputExtension, 1)
	if err := os.MkdirAll(path.Dir(outputFile), 0755); err != nil {
		return "", err
	}
	f, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fmt.Fprintf(f, contents)
	return f.Name(), nil
}
