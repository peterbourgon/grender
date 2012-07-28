package main

import (
	"github.com/hoisie/mustache"
	"github.com/russross/blackfriday"
)

// RenderTemplate renders templateFile (a relative filename underneath
// templatePath) using the given Context, returning a buffer.
// RenderTemplate will also compose any included templates.
func RenderTemplate(templatePath, templateFile string, ctx Context) ([]byte, error) {
	buf, err := ComposeTemplate(templatePath, templateFile)
	if err != nil {
		return nil, err
	}
	rendered := mustache.Render(string(buf), ctx)
	return []byte(rendered), nil
}

// RenderMarkdown renders the passed buffer as Markdown, supressing errors.
func RenderMarkdown(buf []byte) []byte {
	htmlOptions := 0
	htmlOptions = htmlOptions | blackfriday.HTML_GITHUB_BLOCKCODE
	htmlOptions = htmlOptions | blackfriday.HTML_USE_SMARTYPANTS
	htmlRenderer := blackfriday.HtmlRenderer(
		htmlOptions,
		"", // title
		"", // css
	)
	mdOptions := 0
	mdOptions = mdOptions | blackfriday.EXTENSION_FENCED_CODE
	out := blackfriday.Markdown(
		buf,
		htmlRenderer,
		mdOptions,
	)
	return out
}
