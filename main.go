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
	verbose   = flag.Bool("verbose", false, "print verbose information")
	sourceDir = flag.String("source", "src", "path to site source (input)")
	targetDir = flag.String("target", "tgt", "path to site target (output)")
	globalKey = flag.String("global.key", "files", "template node name for per-file metadata")
)

func init() {
	flag.Parse()

	if *debug {
		*verbose = true
	}

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
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		switch filepath.Ext(path) {
		case ".json":
			metadata := mustJSON(mustRead(path))
			s.Add(filepath.Dir(path), metadata)
			Infof("%s gathered (%d element(s))", path, len(metadata))
		}
		return nil
	}
}

func gatherSource(s StackReadWriter, m map[string]interface{}) filepath.WalkFunc {
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		switch filepath.Ext(path) {
		case ".html":
			fullMetadata := map[string]interface{}{
				"source":  diffPath(*sourceDir, path),
				"target":  diffPath(*targetDir, targetFor(path, filepath.Ext(path))),
				"url":     "/" + diffPath(*targetDir, targetFor(path, filepath.Ext(path))),
				"sortkey": filepath.Base(path),
			}
			metadataBuf, _ := splitMetadata(mustRead(path))
			if len(metadataBuf) > 0 {
				fileMetadata := mustJSON(metadataBuf)
				s.Add(path, fileMetadata)
			}
			fullMetadata = mergemap.Merge(fullMetadata, s.Get(path))
			splatInto(m, diffPath(*sourceDir, path), fullMetadata)
			Infof("%s gathered (%d element(s))", path, len(fullMetadata))

		case ".md":
			fullMetadata := map[string]interface{}{
				"source":  diffPath(*sourceDir, path),
				"target":  diffPath(*targetDir, targetFor(path, ".html")),
				"url":     "/" + diffPath(*targetDir, targetFor(path, ".html")),
				"sortkey": filepath.Base(path),
			}
			metadataBuf, _ := splitMetadata(mustRead(path))
			if len(metadataBuf) > 0 {
				fileMetadata := mustJSON(metadataBuf)
				s.Add(path, fileMetadata)
			}
			fullMetadata = mergemap.Merge(fullMetadata, s.Get(path))
			splatInto(m, diffPath(*sourceDir, path), fullMetadata)
			Infof("%s gathered (%d element(s))", path, len(fullMetadata))
		}
		return nil
	}
}

func transform(s StackReader) filepath.WalkFunc {
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		switch filepath.Ext(path) {
		case ".json":
			Infof("%s ignored for transformation", path)

		case ".html":
			// read
			_, contentBuf := splitMetadata(mustRead(path))

			// render
			metadata := s.Get(path)
			outputBuf := renderTemplate(path, contentBuf, metadata)

			// write
			dst := targetFor(path, filepath.Ext(path))
			mustWrite(dst, outputBuf)
			Infof("%s transformed to %s", path, dst)

		case ".md":
			// read
			_, contentBuf := splitMetadata(mustRead(path))

			// render the markdown, and put it into the 'content' key of an
			// interstitial metadata, to be fed to the template renderer
			metadata := mergemap.Merge(s.Get(path), map[string]interface{}{
				"content": template.HTML(renderMarkdown(contentBuf)),
			})

			// render the complete html output according to the template
			outputBuf := renderTemplate(path, mustTemplate(s, path), metadata)

			// write
			dst := targetFor(path, ".html")
			mustWrite(dst, outputBuf)
			Infof("%s transformed to %s", path, dst)

		case ".source", ".template":
			Infof("%s ignored for transformation", path)

		default:
			dst := targetFor(path, filepath.Ext(path))
			mustCopy(dst, path)
			Infof("%s transformed to %s verbatim", path, dst)
		}
		return nil
	}
}

func renderTemplate(path string, input []byte, metadata map[string]interface{}) []byte {
	funcMap := template.FuncMap{
		"importcss": func(filename string) template.CSS {
			filename = filepath.Join(filepath.Dir(path), filename)
			return template.CSS(mustRead(filename))
		},
		"importjs": func(filename string) template.JS {
			filename = filepath.Join(filepath.Dir(path), filename)
			return template.JS(mustRead(filename))
		},
		"importhtml": func(filename string) template.HTML {
			filename = filepath.Join(filepath.Dir(path), filename)
			return template.HTML(mustRead(filename))
		},
		"sorted": sortedValues,
	}
	templateName := diffPath(*sourceDir, path)
	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(string(input))
	if err != nil {
		Fatalf("%s", err)
	}
	output := bytes.Buffer{}
	if err := tmpl.Execute(&output, metadata); err != nil {
		Fatalf("%s", err)
	}
	return output.Bytes()
}

func renderMarkdown(input []byte) []byte {
	htmlOptions := 0
	htmlOptions = htmlOptions | blackfriday.HTML_GITHUB_BLOCKCODE
	htmlOptions = htmlOptions | blackfriday.HTML_USE_SMARTYPANTS
	title, css := "", ""
	htmlRenderer := blackfriday.HtmlRenderer(htmlOptions, title, css)

	mdOptions := 0
	mdOptions = mdOptions | blackfriday.EXTENSION_FENCED_CODE

	return blackfriday.Markdown(input, htmlRenderer, mdOptions)
}
