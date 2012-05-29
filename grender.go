package main

import (
	"flag"
	"github.com/peterbourgon/bonus/xlog"
	"github.com/hoisie/mustache"
	"strings"
	"os"
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
		ctx := GetContext(*sourceDir, page)
		templateFile, err := GetTemplate(*sourceDir, *templatesDir, page)
		if err != nil {
			xlog.Problemf("%s: %s", page, err)
			continue
		}
		tmpl := mustache.RenderFile(templateFile, ctx)
		// TODO make a function
		outputFile := strings.Replace(page, *sourceDir, *outputDir, 1)
		outputFile = strings.Replace(outputFile, PageExtension, OutputExtension, 1)
		f, err := os.Create(outputFile)
		if err != nil {
			xlog.Problemf("%s: %s", outputFile, err)
			continue
		}
		fmt.Fprintf(f, tmpl)
		xlog.Infof("%s: wrote %d bytes to %s", page, len(tmpl), outputFile)
		f.Close()
	}
}
