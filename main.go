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
	debug     = flag.Bool("debug", false, "print debug information (implies verbose)")
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
	filepath.Walk(*sourceDir, gatherJSON(s))
	filepath.Walk(*sourceDir, gatherSource(s, m))
	s.Add("", map[string]interface{}{*globalKey: m})
	filepath.Walk(*sourceDir, transform(s))
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

func gatherJSON(s StackReadWriter) filepath.WalkFunc {
	Debugf("gathering JSON")
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		switch filepath.Ext(path) {
		case ".json":
			metadata := ParseJSON(Read(path))
			s.Add(filepath.Dir(path), metadata)
			Infof("%s gathered (%d element(s))", path, len(metadata))
		}
		return nil
	}
}

func gatherSource(s StackReadWriter, m map[string]interface{}) filepath.WalkFunc {
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
			Infof("%s gathered (%d element(s))", path, len(fullMetadata))

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
			Infof("%s gathered (%d element(s))", path, len(fullMetadata))
		}
		return nil
	}
}

func transform(s StackReader) filepath.WalkFunc {
	Debugf("transforming")
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			Debugf("descending into %s", path)
			return nil // descend
		}

		Debugf("processing %s", path)
		switch filepath.Ext(path) {
		case ".json":
			Infof("%s ignored for transformation", path)

		case ".html":
			// read
			_, contentBuf := splitMetadata(Read(path))

			// render
			metadata := mergemap.Merge(s.Get(path), map[string]interface{}{
				"content": template.HTML(contentBuf),
			})
			var outputBuf []byte
			if templateBuf, err := MaybeTemplate(s, path); err == nil {
				outputBuf = renderTemplate(path, templateBuf, metadata)
			} else {
				outputBuf = renderTemplate(path, contentBuf, metadata)
			}

			// write
			dst := TargetFor(path, filepath.Ext(path))
			Write(dst, outputBuf)
			Infof("%s transformed to %s", path, dst)

		case ".md":
			// read
			_, contentBuf := splitMetadata(Read(path))

			// render
			metadata := mergemap.Merge(s.Get(path), map[string]interface{}{
				"content": template.HTML(renderMarkdown(contentBuf)),
			})
			template := Template(s, path)
			outputBuf := renderTemplate(path, template, metadata)

			// write
			dst := TargetFor(path, ".html")
			Write(dst, outputBuf)
			Infof("%s transformed to %s", path, dst)

		case ".source", ".template":
			Infof("%s ignored for transformation", path)

		default:
			dst := TargetFor(path, filepath.Ext(path))
			Copy(dst, path)
			Infof("%s transformed to %s verbatim", path, dst)
		}
		return nil
	}
}

func renderTemplate(path string, input []byte, metadata map[string]interface{}) []byte {
	Debugf("rendering template (%s) with %d byte(s) of input", path, len(input))
	funcMap := template.FuncMap{
		"importcss": func(filename string) template.CSS {
			Debugf("importcss: importing %s from %s", filename, path)
			filename = filepath.Join(*sourceDir, filename)
			Debugf("importcss: looking for %s", filename)
			return template.CSS(Read(filename))
		},
		"importjs": func(filename string) template.JS {
			Debugf("importjs: importing %s from %s", filename, path)
			filename = filepath.Join(*sourceDir, filename)
			Debugf("importjs: looking for %s", filename)
			return template.JS(Read(filename))
		},
		"importhtml": func(filename string) template.HTML {
			Debugf("importhtml: importing %s from %s", filename, path)
			filename = filepath.Join(*sourceDir, filename)
			Debugf("importhtml: looking for %s", filename)
			return template.HTML(Read(filename))
		},
		"sorted": SortedValues,
	}

	templateName := Relative(*sourceDir, path)
	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(string(input))
	if err != nil {
		Fatalf("render template: parsing: %s", err)
	}

	output := bytes.Buffer{}
	if err := tmpl.Execute(&output, metadata); err != nil {
		Fatalf("render template: executing: %s", err)
	}

	return output.Bytes()
}

func renderMarkdown(input []byte) []byte {
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
