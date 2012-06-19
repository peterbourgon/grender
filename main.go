package main

import (
	"flag"
	"fmt"
	"github.com/peterbourgon/bonus/xlog"
)

var (
	debug       *bool   = flag.Bool("debug", false, "enable debug output")
	templateDir *string = flag.String("template-dir", "_templates", "templates directory")
	sourceDir   *string = flag.String("source-dir", "_source", "source directory")
	staticDir   *string = flag.String("static-dir", "_static", "static directory")
)

type Context map[string]interface{}

func init() {
	flag.Parse()
	xlog.Initialize(*debug)
}

// A source file will specify
//  - the template it wants to use to render itself
//  - the data it wants to put into that template
//  - its output path
//
// A template file will specify
//  - subtemplates it can pull in
//  - somehow defaults for unspecified variables

func main() {
	ctx := map[string]interface{}{
		"title": "Contextual title",
	}
	buf, err := Render(*templateDir, "index.tmpl", ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(buf))
}
