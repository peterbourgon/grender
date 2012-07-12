package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/peterbourgon/bonus/xlog"
	"strings"
)

var (
	debug       *bool   = flag.Bool("debug", false, "enable debug output")
	templateDir *string = flag.String("template-dir", "_templates", "templates directory")
	sourceDir   *string = flag.String("source-dir", "_source", "source directory")
	staticDir   *string = flag.String("static-dir", "_static", "static directory")
	outputDir   *string = flag.String("output-dir", "_output", "output directory")
)

func init() {
	flag.Parse()
	xlog.Initialize(*debug)
}

// A source file will specify
//  - the template it wants to use to render itself
//  - the data it wants to put into that template
//  - its output path
//
// A template file will specify
//  - subtemplates it can pull in
//  - somehow defaults for unspecified variables

func main() {
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		xlog.Fatalf("%s: %s", *outputDir, err)
	}

	// copy static
	copyStatic := exec.Command("cp", "-r", *staticDir+"/", *outputDir+"/")
	buf, err := copyStatic.CombinedOutput()
	if err != nil {
		xlog.Fatalf("copying %s: %s (%s)", *staticDir, err, strings.TrimSpace(string(buf)))
	}

	// render source files
	sourceFiles, err := filepath.Glob(*sourceDir + "/*")
	if err != nil {
		xlog.Fatalf("%s: %s", *sourceDir, err)
	}
	for _, sourceFile := range sourceFiles {
		// extract context
		templateFile, ctx, outputFile, err := ContextFrom(sourceFile)
		if err != nil {
			xlog.Problemf("Context: %s: %s", sourceFile, err)
			continue
		}

		// render using specified template
		buf, err := RenderTemplate(*templateDir, templateFile, ctx)
		if err != nil {
			xlog.Problemf("Render: %s: %s", sourceFile, err)
			continue
		}

		// write output
		totalOutputFile := *outputDir + "/" + outputFile
		if err := os.MkdirAll(filepath.Dir(totalOutputFile), 0755); err != nil {
			xlog.Problemf("MkdirAll: %s: %s: %s", sourceFile, outputFile, err)
			continue
		}
		f, err := os.Create(totalOutputFile)
		if err != nil {
			xlog.Problemf("Create: %s: %s: %s", sourceFile, outputFile, err)
			continue
		}
		defer f.Close()
		n, err := f.Write(buf)
		if err != nil {
			xlog.Problemf("Write: %s: %s: %s", sourceFile, outputFile, err)
			continue
		}
		if n != len(buf) {
			xlog.Problemf("%s: %s: %d < %d", sourceFile, outputFile, n, len(buf))
			continue
		}
	}
}
