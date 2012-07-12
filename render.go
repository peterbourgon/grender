package main

import (
	"github.com/russross/blackfriday"
	"github.com/hoisie/mustache"
)

const (
	MarkdownKey = "markdown"
)

func RenderTemplate(templateDir, inFile string, ctx Context) ([]byte, error) {
	// First, expand all Composition directives.
	buf, err := ComposeTemplate(templateDir, inFile)
	if err != nil {
		return nil, err
	}

	// Second, inject MarkdownKey into the ctx, from BodyKey content (if any)
	if i, ok := ctx[BodyKey]; ok {
		if s, ok := i.(string); ok {
			ctx[MarkdownKey] = string(blackfriday.MarkdownBasic([]byte(s)))
		}
	}

	// Last, render the template
	rendered := mustache.Render(string(buf), ctx)
	return []byte(rendered), nil
}
