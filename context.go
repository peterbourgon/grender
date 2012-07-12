package main

import (
	// "github.com/kylelemons/go-gypsy/yaml"
	"launchpad.net/goyaml"
	"io/ioutil"
	"os"
	"strings"
)

const (
	SourceSeparator = "---"
	BodyKey         = "body"
)

type Context map[string]interface{}

func ContextFrom(sourceFile string) (Context, error) {
	ctx := Context{}

	// open and read the sourceFile
	f, err := os.Open(sourceFile)
	if err != nil {
		return ctx, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return ctx, err
	}
	s := string(buf)

	// split on SourceSeparator; treat second part as 'body'
	if idx := strings.Index(s, SourceSeparator); idx >= 0 {
		ctx[BodyKey] = string(buf[idx+len(SourceSeparator)+1:])
		buf = buf[:idx]
	}

	err = goyaml.Unmarshal(buf, ctx)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}
