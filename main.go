package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"log"
)

var (
	debug       = flag.Bool("debug", false, "enable debug output")
	templateDir = flag.String("template-dir", "templates", "templates directory")
	postsDir    = flag.String("posts-dir", "posts", "posts directory")
	sourceDir   = flag.String("source-dir", "source", "other source files directory")
	staticDir   = flag.String("static-dir", "static", "static directory")
	outputDir   = flag.String("output-dir", "site", "output directory")
	postsSubdir = flag.String("posts-subdir", "posts", "posts subdirectory (in output)")
)

func init() {
	flag.Parse()
	log.SetFlags(0)
}

const (
	TemplateKey         = "template"
	PostsKey            = "posts"
	TitleKey            = "title"
	OutputDirectoryMode = 0755
)

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

	// copy static files
	for _, path := range FilesIn(*staticDir) {
		d := *staticDir + "/" + filepath.Dir(path)
		if err := os.MkdirAll(d, OutputDirectoryMode); err != nil {
			Problemf("%s: %s", d, err)
			continue
		}
		tgt := *staticDir + "/" + path
		b, err := exec.Command("cp", path, tgt).CombinedOutput()
		if err != nil {
			Fatalf("cp %s %s: %s (%s)", path, tgt, err, strings.TrimSpace(string(b)))
		}
	}

	// render and copy posts
	posts := PostsIn(*postsDir, "")
	for i, post := range posts {
		ctx, err := ContextFrom(post.Filename)

		// extract title from context
		if o, ok := ctx[TitleKey]; ok {
			if s, ok := o.(string); ok {
				posts[i].Title = s
			}
		}

		// extract template to use from context
		o, ok := ctx[TemplateKey]
		if !ok {
			Problemf("Extract: %s: no %s", post.Filename, TemplateKey)
			continue
		}
		templateFile, ok := o.(string)
		if !ok {
			Problemf("Extract: %s: %s: bad type", post.Filename, TemplateKey)
			continue
		}
		buf, err := RenderTemplate(*templateDir, templateFile, ctx)
		if err != nil {
			Problemf("Render: %s: %s", post.Filename, err)
			continue
		}

		outputFile := *outputDir + "/" + *postsSubdir
		outputFile += fmt.Sprintf("/%s.html", filepath.Base(post.Filename))
		d := filepath.Dir(outputFile)
		if err := os.MkdirAll(d, OutputDirectoryMode); err != nil {
			Problemf("MkdirAll: %s: %s: %s", post.Filename, outputFile, err)
			continue
		}
		f, err := os.Create(outputFile)
		if err != nil {
			Problemf("Create: %s: %s: %s", post.Filename, outputFile, err)
			continue
		}
		defer f.Close()
		n, err := f.Write(buf)
		if err != nil {
			Problemf("Write: %s: %s: %s", post.Filename, outputFile, err)
			continue
		}
		if n != len(buf) {
			Problemf("Write: %s: %s: %d < %d", post.Filename, outputFile, n, len(buf))
			continue
		}
		log.Printf("%s -> %s OK", post.Filename, outputFile)
	}

	// render and copy source files
	for _, sourceFile := range FilesIn(*sourceDir) {
		ctx, err := ContextFrom(sourceFile)
		ctx[PostsKey] = posts

		i, ok := ctx[TemplateKey]
		if !ok {
			Problemf("Extract: %s: no %s", sourceFile, TemplateKey)
			continue
		}
		templateFile, ok := i.(string)
		if !ok {
			Problemf("Extract: %s: %s: bad type", sourceFile, TemplateKey)
			continue
		}
		buf, err := RenderTemplate(*templateDir, templateFile, ctx)
		if err != nil {
			Problemf("Render: %s: %s", sourceFile, err)
			continue
		}

		outputFile := fmt.Sprintf("%s/%s.html", *outputDir, filepath.Base(sourceFile))
		d := filepath.Dir(outputFile)
		if err := os.MkdirAll(d, OutputDirectoryMode); err != nil {
			Problemf("MkdirAll: %s: %s: %s", sourceFile, outputFile, err)
			continue
		}
		f, err := os.Create(outputFile)
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
			Problemf("Write: %s: %s: %d < %d", sourceFile, outputFile, n, len(buf))
			continue
		}
		log.Printf("%s -> %s OK", sourceFile, outputFile)
	}

	/*

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
	*/
}
