package main

import (
	// "github.com/kylelemons/go-gypsy/yaml"
	"launchpad.net/goyaml"
	"io/ioutil"
	"fmt"
	"os"
	"strings"
)

const (
	SourceSeparator = "---"
	BodyKey         = "body"
	TemplateKey     = "template"
	OutputKey       = "output"
)

type Context map[string]interface{}

func ContextFrom(sourceFile string) (string, Context, string, error) {
	t, ctx, o := "", Context{}, ""

	// open and read the sourceFile
	f, err := os.Open(sourceFile)
	if err != nil {
		return t, ctx, o, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return t, ctx, o, err
	}
	s := string(buf)

	// split on SourceSeparator; treat second part as 'body'
	if idx := strings.Index(s, SourceSeparator); idx >= 0 {
		ctx[BodyKey] = string(buf[idx+len(SourceSeparator)+1:])
		buf = buf[:idx]
	}

	err = goyaml.Unmarshal(buf, ctx)
	if err != nil {
		return t, ctx, o, err
	}

	// extract necessary fields from sourceFile
	for _, pair := range []struct {
		key    string
		target *string
	}{
		{TemplateKey, &t},
		{OutputKey, &o},
	} {
		var i interface{}
		var ok bool
		i, ok = ctx[pair.key]
		if !ok {
			return t, ctx, o, fmt.Errorf("no '%s'", pair.key)
		}
		*(pair.target), ok = i.(string)
		if !ok {
			return t, ctx, o, fmt.Errorf("invalid '%s' type", pair.key)
		}
		delete(ctx, pair.key)
	}

	// i, ok = ctx[TemplateKey]
	// if !ok {
	// 	return t, ctx, o, fmt.Errorf("no '%s'", TemplateKey)
	// }
	// t, ok = i.(string)
	// if !ok {
	// 	return t, ctx, o, fmt.Errorf("invalid '%s' type", TemplateKey)
	// }
	// delete(ctx, TemplateKey)

	// i, ok = ctx[OutputKey]
	// if !ok {
	// 	return t, ctx, o, fmt.Errorf("no '%s'", OutputKey)
	// }
	// o, ok = i.(string)
	// if !ok {
	// 	return t, ctx, o, fmt.Errorf("invalid '%s' type", OutputKey)
	// }
	// delete(ctx, OutputKey)

	// done
	return t, ctx, o, nil
}
