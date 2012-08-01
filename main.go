package main

import (
	"flag"
	"fmt"
)

var (
	metadataDelimiter = flag.String("metadata-delimiter", "---", "string which terminates the metadata field in source files")
	contentKey        = flag.String("content-key", "content", "the metadata key in which the body of a source file is provided in contexts")
	templateKey       = flag.String("template-key", "template", "the metadata key that specifies the desired template")
	outputKey         = flag.String("output-key", "output", "the metadata key that specifies the desired output file basename")
	indexKey          = flag.String("index-key", "index", "the metadata key that specifies the global index")
	defaultIndexType  = flag.String("default-index-type", "blog", "the default 'type' of index metadata, if none is explicitly provided")
	indexContentCount = flag.Int("index-content-count", 1, "the first N index-tuples whose content will be provided in each typed index")
	outputExtension   = flag.String("output-extension", "html", "all rendered output files will have this extension")
	sourcePath        = flag.String("source-path", "_source", "where source files are contained")
	templatePath      = flag.String("template-path", "_templates", "where template files are contained")
	staticPath        = flag.String("static-path", "_static", "static files copied directly to -output-path")
	outputPath        = flag.String("output-path", "_site", "where grender will place output files")
	blogPath          = flag.String("blog-path", "blog", "path under -output-path where grender will place blog entries")
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
	idx := NewIndex()
	contexts := map[string]Context{}
	for _, sourceFile := range Filenames(*sourcePath) {
		ctx, err := ParseSourceFile(
			*sourcePath,
			sourceFile,
			idx,
			*blogPath,
			*metadataDelimiter,
			*contentKey,
			*templateKey,
			*indexKey,
			*outputKey,
			*outputExtension,
		)
		if err != nil {
			Logf("%s: parsing: %s", sourceFile, err)
			continue
		}
		contexts[sourceFile] = ctx
	}

	// TODO manipulate index

	// Second pass: render source files
	for sourceFile, ctx := range contexts {
		templateFile, err := ctx.GetString(*templateKey)
		if err != nil {
			Logf("%s: parsing: %s", sourceFile, err)
			continue
		}

		outputFile, err := ctx.GetString(*outputKey)
		if err != nil {
			Logf("%s: parsing: %s", sourceFile, err)
			continue
		}

		ctx[*indexKey] = idx
		output, err := RenderTemplate(
			*templatePath,
			templateFile,
			ctx,
		)
		if err != nil {
			Logf("%s: rendering: %s", sourceFile, err)
			continue
		}

		outputFile = *outputPath + "/" + outputFile + "." + *outputExtension
		if err := WriteOutput(output, outputFile); err != nil {
			Logf("%s: writing: %s", sourceFile, err)
			continue
		}
	}
}
