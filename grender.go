package main

import (
	"flag"
	"github.com/peterbourgon/bonus/xlog"
	// "github.com/hoisie/mustache"
	// "fmt"
)

var (
	debug		*bool	= flag.Bool("debug", false, "enable debug logging")
	templatesDir	*string	= flag.String("templates-dir", "_templates", "directory containing templates")
	sourceDir	*string	= flag.String("source-dir", "_source", "directory containing source")
	outputDir	*string	= flag.String("output-dir", "_output", "directory where site will be written")
	temp		*string	= flag.String("temp", "", "temporary flag")
)

func init() {
	flag.Parse()
	xlog.Initialize(*debug)
}

func main() {
	xlog.Debugf("starting")
}
