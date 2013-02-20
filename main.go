package main

import (
	"bytes"
	"flag"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

var (
	FrontSeparator = []byte("---\n")
)

var (
	sourceDir = flag.String("source", "src", "path to site source (input)")
	targetDir = flag.String("target", "tgt", "path to site target (output)")
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
}

func main() {
	filepath.Walk(*sourceDir, sourceWalk())
}

func sourceWalk() filepath.WalkFunc {
	s := NewStack()

	readAndAdd := func(path string) {
		m := mustJSON(mustRead(path))
		s.Add(filepath.Dir(path), m)
		log.Printf("%-50s added to stack: %v", path, m)
	}

	cp := func(path string) {
		dst := filepath.Join(*targetDir, diffPath(*sourceDir, path))
		copyFile(dst, path)
		log.Printf("%-50s copied to %s", path, dst)
	}

	type specificRenderFunc func(input []byte, m map[string]interface{}) []byte

	frontMatter := func(f specificRenderFunc) func(path string) {
		return func(path string) {
			input := mustRead(path)
			split := bytes.SplitN(input, FrontSeparator, 2)
			if len(split) == 2 {
				s.Add(path, mustJSON(split[0]))
				input = split[1]
			}
			output := f(input, s.Get(path))
			dst := targetFor(path)
			mustWrite(dst, output)
			log.Printf("%-50s rendered to %s", path, dst)
		}
	}

	renderHTML := func(input []byte, m map[string]interface{}) []byte {
		funcMap := template.FuncMap{
			"importcss": func(filename string) template.CSS {
				return template.CSS(mustRead(filepath.Join(*sourceDir, filename)))
			},
			"importjs": func(filename string) template.JS {
				return template.JS(mustRead(filepath.Join(*sourceDir, filename)))
			},
			"importhtml": func(filename string) template.HTML {
				return template.HTML(mustRead(filepath.Join(*sourceDir, filename)))
			},
		}
		tmpl, err := template.New("x").Funcs(funcMap).Parse(string(input))
		if err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}
		output := bytes.Buffer{}
		if err := tmpl.Execute(&output, m); err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}
		return output.Bytes()
	}

	renderMarkdown := func(input []byte, m map[string]interface{}) []byte {
		input = renderHTML(input, m)

		htmlOptions := 0
		htmlOptions = htmlOptions | blackfriday.HTML_GITHUB_BLOCKCODE
		htmlOptions = htmlOptions | blackfriday.HTML_USE_SMARTYPANTS
		title, css := "", ""
		htmlRenderer := blackfriday.HtmlRenderer(htmlOptions, title, css)

		mdOptions := 0
		mdOptions = mdOptions | blackfriday.EXTENSION_FENCED_CODE

		return blackfriday.Markdown(input, htmlRenderer, mdOptions)
	}

	noop := func(string) {}

	ext := map[string]func(path string){
		".json":   readAndAdd,
		".html":   frontMatter(renderHTML),
		".md":     frontMatter(renderMarkdown),
		".source": noop,
	}

	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}
		f, ok := ext[filepath.Ext(path)]
		if !ok {
			f = cp
		}
		f(path)
		return nil
	}
}
