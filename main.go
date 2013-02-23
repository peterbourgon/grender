package main

import (
	"bytes"
	"flag"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	FrontSeparator = []byte("---\n")
)

var (
	sourceDir  = flag.String("source", "src", "path to site source (input)")
	targetDir  = flag.String("target", "tgt", "path to site target (output)")
	globalKeys = flag.String("global", "blog", "comma-separated list of global keys")
)

var (
	Globals = map[string]struct{}{}
)

func init() {
	log.SetFlags(0)
	flag.Parse()

	var err error
	for _, s := range []*string{sourceDir, targetDir} {
		if *s, err = filepath.Abs(*s); err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}
	}

	for _, k := range strings.Split(*globalKeys, ",") {
		Globals[k] = struct{}{}
	}
}

func main() {
	s := NewStack()
	filepath.Walk(*sourceDir, gatherGlobals(s))
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

func gatherGlobals(s *Stack) filepath.WalkFunc {
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		switch filepath.Ext(path) {
		case ".json":
			for k, v := range mustJSON(mustRead(path)) {
				if _, ok := Globals[k]; ok {
					subMetadata := map[string]interface{}{k: v}
					s.Add(filepath.Dir(path), subMetadata) // this dir
				}
			}
		case ".html", ".md":
			metadataBuf, _ := splitMetadata(mustRead(path))
			if len(metadataBuf) > 0 {
				for k, v := range mustJSON(metadataBuf) {
					if _, ok := Globals[k]; ok {
						subMetadata := map[string]interface{}{k: v}
						s.Add(path, subMetadata) // this file only
					}
				}
			}
		}
		return nil
	}
}

func renderTemplate(path string, input []byte, metadata map[string]interface{}) []byte {
	funcMap := template.FuncMap{
		"importcss": func(filename string) template.CSS {
			return template.CSS(mustRead(filepath.Join(filepath.Dir(path), filename)))
		},
		"importjs": func(filename string) template.JS {
			return template.JS(mustRead(filepath.Join(filepath.Dir(path), filename)))
		},
		"importhtml": func(filename string) template.HTML {
			return template.HTML(mustRead(filepath.Join(filepath.Dir(path), filename)))
		},
	}
	tmpl, err := template.New("x").Funcs(funcMap).Parse(string(input))
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	output := bytes.Buffer{}
	if err := tmpl.Execute(&output, metadata); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
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

func transform(s *Stack) filepath.WalkFunc {
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		switch filepath.Ext(path) {
		case ".json":
			s.Add(filepath.Dir(path), mustJSON(mustRead(path)))
			log.Printf("%s added to stack", path)

		case ".html":
			metadataBuf, dataBuf := splitMetadata(mustRead(path))
			if len(metadataBuf) > 0 {
				s.Add(path, mustJSON(metadataBuf))
			}
			outputBuf := renderTemplate(path, dataBuf, s.Get(path))
			dst := targetFor(path, filepath.Ext(path))
			mustWrite(dst, outputBuf)
			log.Printf("%s rendered to %s", path, dst)

		case ".md":
			// .md may contain front matter
			metadataBuf, dataBuf := splitMetadata(mustRead(path))
			if len(metadataBuf) > 0 {
				s.Add(path, mustJSON(metadataBuf))
			}

			// render the markdown, and put it into the 'content' key of an
			// interstitial metadata, to be fed to the template renderer
			myMetadata := mergeInto(s.Get(path), map[string]interface{}{
				"content": renderMarkdown(dataBuf),
			})

			// render the complete html output according to the template
			outputBuf := renderTemplate(path, mustTemplate(s, path), myMetadata)

			// write
			dst := targetFor(path, ".html") // force .html extension
			mustWrite(dst, outputBuf)
			log.Printf("%s rendered to %s", path, dst)

		case ".source", ".template":
			log.Printf("%s ignored", path)

		default:
			dst := filepath.Join(*targetDir, diffPath(*sourceDir, path))
			mustCopy(targetFor(path, filepath.Ext(path)), path)
			log.Printf("%s copied to %s verbatim", path, dst)
		}
		return nil
	}
}
