package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"log"
)

var (
	debug       = flag.Bool("debug", false, "enable debug output")
	templateDir = flag.String("template-dir", "_templates", "templates directory")
	sourceDir   = flag.String("source-dir", "_source", "source directory")
	staticDir   = flag.String("static-dir", "_static", "static directory")
	outputDir   = flag.String("output-dir", "_output", "output directory")
)

func Debugf(format string, args ...interface{}) {
	if *debug {
		log.Printf("DEBUG "+format, args...)
	}
}

func Problemf(format string, args ...interface{}) {
	log.Printf("PROBLEM "+format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf("FATAL "+format, args...)
}

func init() {
	flag.Parse()
	log.SetFlags(0)
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
		Fatalf("%s: %s", *outputDir, err)
	}

	// copy static
	copyStatic := exec.Command("cp", "-r", *staticDir+"/", *outputDir+"/")
	buf, err := copyStatic.CombinedOutput()
	if err != nil {
		Fatalf("copying %s: %s (%s)", *staticDir, err, strings.TrimSpace(string(buf)))
	}

	// render source files
	sourceFiles, err := filepath.Glob(*sourceDir + "/*")
	if err != nil {
		Fatalf("%s: %s", *sourceDir, err)
	}
	for i, sourceFile := range sourceFiles {
		// extract context
		templateFile, ctx, outputFile, err := ContextFrom(sourceFile)
		if err != nil {
			Problemf("Context: %s: %s", sourceFile, err)
			continue
		}

		// render using specified template
		buf, err := RenderTemplate(*templateDir, templateFile, ctx)
		if err != nil {
			Problemf("Render: %s: %s", sourceFile, err)
			continue
		}

		// write output
		totalOutputFile := *outputDir + "/" + outputFile
		if err := os.MkdirAll(filepath.Dir(totalOutputFile), 0755); err != nil {
			Problemf("MkdirAll: %s: %s: %s", sourceFile, outputFile, err)
			continue
		}
		f, err := os.Create(totalOutputFile)
		if err != nil {
			Problemf("Create: %s: %s: %s", sourceFile, outputFile, err)
			continue
		}
		defer f.Close()
		n, err := f.Write(buf)
		if err != nil {
			Problemf("Write: %s: %s: %s", sourceFile, outputFile, err)
			continue
		}
		if n != len(buf) {
			Problemf("%s: %s: %d < %d", sourceFile, outputFile, n, len(buf))
			continue
		}

		log.Printf(
			" %d/%d) %s → %s OK",
			i+1,
			len(sourceFiles),
			sourceFile,
			totalOutputFile,
		)
	}
}
