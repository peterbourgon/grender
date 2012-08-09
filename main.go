package main

import (
	"flag"
	"fmt"
)

var (
	sourcePath   = flag.String("source-path", "_source", "where source files are contained")
	templatePath = flag.String("template-path", "_templates", "where template files are contained")
	staticPath   = flag.String("static-path", "_static", "static files copied directly to -output-path")
	outputPath   = flag.String("output-path", "_site", "where grender will place output files")

	blogPath          = flag.String("blog-path", "blog", "path under -output-path where grender will place blog entries")
	outputExtension   = flag.String("output-extension", "html", "all rendered output files will have this extension")
	metadataDelimiter = flag.String("metadata-delimiter", "---", "string which terminates the metadata field in source files")

	contentKey  = flag.String("content-key", "content", "the metadata key in which the body of a source file is provided")
	templateKey = flag.String("template-key", "template", "the metadata key that specifies the desired template")
	outputKey   = flag.String("output-key", "output", "the metadata key that specifies the desired output file basename")
	indexKey    = flag.String("index-key", "index", "the metadata key that specifies the global index in the context")

	debug = flag.Bool("debug", false, "provide debug output")
)

func init() {
	flag.Parse()
}

func Logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Debugf(format string, args ...interface{}) {
	if *debug {
		fmt.Printf("DEBUG "+format+"\n", args...)
	}
}

func main() {
	if err := RecursiveCopy(*staticPath, *outputPath); err != nil {
		Debugf("copying %s to %s: %s", *staticPath, *outputPath, err)
	}

	// First pass: get contexts and build index
	idx, sourceFiles := Index{}, SourceFiles{}
	for _, sourceFile := range Filenames(*sourcePath) {
		sf, err := ParseSourceFile(sourceFile)
		if err != nil {
			Debugf("%s: parsing: %s", sourceFile, err)
			continue
		}
		if sf.Indexable() {
			idx.Add(sf)
		}
		sourceFiles = append(sourceFiles, sf)
	}

	// Second pass: render source files
	for _, sf := range sourceFiles {
		if sf.getBool("index") {
			sf.Metadata["index"] = idx.Render()
		}
		Debugf("%s: rendering with ctx: %v", sf.SourceFile, sf.Metadata)
		output, err := RenderTemplate(
			sf.getString(*templateKey),
			sf.Metadata,
		)
		if err != nil {
			Logf("%s: rendering: %s", sf.SourceFile, err)
			continue
		}

		outputFile := *outputPath + "/" + sf.Output() + "." + *outputExtension
		if err := WriteOutput(output, outputFile); err != nil {
			Logf("%s: writing: %s", sf.SourceFile, err)
			continue
		}
	}
}
