package main

import (
	"flag"
	"fmt"
)

var (
	metadataDelimiter     = flag.String("metadata-delimiter", "---", "string which terminates the metadata field in source files")
	contentKey            = flag.String("content-key", "content", "the metadata key in which the body of a source file is provided in contexts")
	templateKey           = flag.String("template-key", "template", "the metadata key that specifies the desired template")
	outputKey             = flag.String("output-key", "output", "the metadata key that specifies the desired output file basename")
	indexKey              = flag.String("index-key", "index", "the metadata key that specifies the global index in the context")
	indexTupleKey         = flag.String("index-tuple-key", "index", "the metadata key that specifies an index-tuple in a source file metadata")
	indexTupleTypeKey     = flag.String("index-tuple-type-key", "type", "the metadata key that specifies the type of an index-tuple")
	indexTupleSortKeyKey  = flag.String("index-tuple-sortkey-key", "key", "key of an index-tuple used to perform sorting")
	defaultIndexTupleType = flag.String("default-index-tuple-type", "blog", "the default 'type' of index metadata, if none is explicitly provided")
	outputExtension       = flag.String("output-extension", "html", "all rendered output files will have this extension")
	sourcePath            = flag.String("source-path", "_source", "where source files are contained")
	templatePath          = flag.String("template-path", "_templates", "where template files are contained")
	staticPath            = flag.String("static-path", "_static", "static files copied directly to -output-path")
	outputPath            = flag.String("output-path", "_site", "where grender will place output files")
	blogPath              = flag.String("blog-path", "blog", "path under -output-path where grender will place blog entries")
)

func init() {
	flag.Parse()
}

func Logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func main() {
	if err := RecursiveCopy(*staticPath, *outputPath); err != nil {
		Logf("copying %s to %s: %s", *staticPath, *outputPath, err)
	}

	// First pass: get contexts and build index
	idx := Index{}
	sourceFiles := []*SourceFile{}
	for _, sourceFile := range Filenames(*sourcePath) {
		sf, err := ParseSourceFile(sourceFile)
		if err != nil {
			Logf("%s: parsing: %s", sourceFile, err)
			continue
		}
		sf.IndexTuple.ContributeTo(idx) // sorts as it goes
		sourceFiles = append(sourceFiles, sf)
	}

	// Second pass: render source files
	for _, sf := range sourceFiles {
		sf.Metadata[*indexKey] = idx          // provide global index in ctx
		sf.Metadata[*contentKey] = sf.Content // push content in there too
		output, err := RenderTemplate(
			sf.TemplateFile,
			sf.Metadata,
		)
		if err != nil {
			Logf("%s: rendering: %s", sf.SourceFile, err)
			continue
		}

		outputFile := *outputPath + "/" + sf.OutputFile + "." + *outputExtension
		if err := WriteOutput(output, outputFile); err != nil {
			Logf("%s: writing: %s", sf.SourceFile, err)
			continue
		}
	}
}
