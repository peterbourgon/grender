package main

import (
	"bytes"
	"flag"
	"github.com/peterbourgon/mergemap"
	"github.com/russross/blackfriday"
	"html/template"
	"os"
	"path/filepath"
)

var (
	FrontSeparator = []byte("---\n")
)

var (
	debug     = flag.Bool("debug", false, "print debug information")
	sourceDir = flag.String("source", "src", "path to site source (input)")
	targetDir = flag.String("target", "tgt", "path to site target (output)")
	globalKey = flag.String("global.key", "files", "template node name for per-file metadata")
)

func init() {
	flag.Parse()

	var err error
	for _, s := range []*string{sourceDir, targetDir} {
		if *s, err = filepath.Abs(*s); err != nil {
			Fatalf("%s", err)
		}
	}
}

func main() {
	m := map[string]interface{}{}
	s := NewStack()
	filepath.Walk(*sourceDir, GatherJSON(s))
	filepath.Walk(*sourceDir, GatherSource(s, m))
	s.Add("", map[string]interface{}{*globalKey: m})
	filepath.Walk(*sourceDir, Transform(s))
}

// splitMetadata splits the input buffer on FrontSeparator. It returns a byte-
// slice suitable for unmarshaling into metadata, if it exists, and the
// remainder of the input buffer.
func splitMetadata(buf []byte) ([]byte, []byte) {
	split := bytes.SplitN(buf, FrontSeparator, 2)
	if len(split) == 2 {
		return split[0], split[1]
	}
	return []byte{}, buf
}

func GatherJSON(s StackReadWriter) filepath.WalkFunc {
	Debugf("gathering JSON")
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		switch filepath.Ext(path) {
		case ".json":
			metadata := ParseJSON(Read(path))
			s.Add(filepath.Dir(path), metadata)
			Debugf("%s gathered (%d element(s))", path, len(metadata))
		}
		return nil
	}
}

func GatherSource(s StackReadWriter, m map[string]interface{}) filepath.WalkFunc {
	Debugf("gathering source")
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		switch filepath.Ext(path) {
		case ".html":
			fullMetadata := map[string]interface{}{
				"source":  Relative(*sourceDir, path),
				"target":  Relative(*targetDir, TargetFor(path, filepath.Ext(path))),
				"url":     "/" + Relative(*targetDir, TargetFor(path, filepath.Ext(path))),
				"sortkey": filepath.Base(path),
			}
			metadataBuf, _ := splitMetadata(Read(path))
			if len(metadataBuf) > 0 {
				fileMetadata := ParseJSON(metadataBuf)
				s.Add(path, fileMetadata)
			}
			fullMetadata = mergemap.Merge(fullMetadata, s.Get(path))
			SplatInto(m, Relative(*sourceDir, path), fullMetadata)
			Debugf("%s gathered (%d element(s))", path, len(fullMetadata))

		case ".md":
			fullMetadata := map[string]interface{}{
				"source":  Relative(*sourceDir, path),
				"target":  Relative(*targetDir, TargetFor(path, ".html")),
				"url":     "/" + Relative(*targetDir, TargetFor(path, ".html")),
				"sortkey": filepath.Base(path),
			}
			metadataBuf, _ := splitMetadata(Read(path))
			if len(metadataBuf) > 0 {
				fileMetadata := ParseJSON(metadataBuf)
				s.Add(path, fileMetadata)
			}
			fullMetadata = mergemap.Merge(fullMetadata, s.Get(path))
			SplatInto(m, Relative(*sourceDir, path), fullMetadata)
			Debugf("%s gathered (%d element(s))", path, len(fullMetadata))
		}
		return nil
	}
}

func Transform(s StackReader) filepath.WalkFunc {
	Debugf("transforming")
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			Debugf("descending into %s", path)
			return nil // descend
		}

		Debugf("processing %s", path)
		switch filepath.Ext(path) {
		case ".json":
			Debugf("%s ignored for transformation", path)

		case ".html":
			// read
			_, contentBuf := splitMetadata(Read(path))

			// render
			outputBuf := RenderTemplate(path, contentBuf, s.Get(path))

			// write
			dst := TargetFor(path, filepath.Ext(path))
			Write(dst, outputBuf)
			Debugf("%s transformed to %s", path, dst)

		case ".md":
			// read
			_, contentBuf := splitMetadata(Read(path))

			// render
			metadata := mergemap.Merge(s.Get(path), map[string]interface{}{
				"content": template.HTML(RenderMarkdown(contentBuf)),
			})
			templatePath, templateBuf := Template(s, path)
			outputBuf := RenderTemplate(templatePath, templateBuf, metadata)

			// write
			dst := TargetFor(path, ".html")
			Write(dst, outputBuf)
			Debugf("%s transformed to %s", path, dst)

		case ".source", ".template":
			Debugf("%s ignored for transformation", path)

		default:
			dst := TargetFor(path, filepath.Ext(path))
			Copy(dst, path)
			Debugf("%s transformed to %s verbatim", path, dst)
		}
		return nil
	}
}

func RenderTemplate(path string, input []byte, metadata map[string]interface{}) []byte {
	R := func(relativeFilename string) string {
		filename := filepath.Join(filepath.Dir(path), relativeFilename)
		return string(RenderTemplate(filename, Read(filename), metadata))
	}
	importhtml := func(relativeFilename string) template.HTML {
		return template.HTML(R(relativeFilename))
	}
	importcss := func(relativeFilename string) template.CSS {
		return template.CSS(R(relativeFilename))
	}
	importjs := func(relativeFilename string) template.JS {
		return template.JS(R(relativeFilename))
	}

	templateName := Relative(*sourceDir, path)
	funcMap := template.FuncMap{
		"importhtml": importhtml,
		"importcss":  importcss,
		"importjs":   importjs,
		"sorted":     SortedValues,
	}

	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(string(input))
	if err != nil {
		Fatalf("Render Template %s: Parse: %s", path, err)
	}

	output := bytes.Buffer{}
	if err = tmpl.Execute(&output, metadata); err != nil {
		Fatalf("Render Template %s: Execute: %s", path, err)
	}

	return output.Bytes()
}

func RenderMarkdown(input []byte) []byte {
	Debugf("rendering %d byte(s) of Markdown", len(input))
	htmlOptions := 0
	htmlOptions = htmlOptions | blackfriday.HTML_GITHUB_BLOCKCODE
	htmlOptions = htmlOptions | blackfriday.HTML_USE_SMARTYPANTS
	title, css := "", ""
	htmlRenderer := blackfriday.HtmlRenderer(htmlOptions, title, css)

	mdOptions := 0
	mdOptions = mdOptions | blackfriday.EXTENSION_FENCED_CODE

	return blackfriday.Markdown(input, htmlRenderer, mdOptions)
}
