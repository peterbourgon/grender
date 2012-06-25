package main

import (
	"github.com/hoisie/mustache"
)

func RenderTemplate(templateDir, inFile string, ctx Context) ([]byte, error) {
	// First, expand all Composition directives.
	buf, err := ComposeTemplate(templateDir, inFile)
	if err != nil {
		return nil, err
	}

	// Then, render the template
	rendered := mustache.Render(string(buf), ctx)
	return []byte(rendered), nil
}
