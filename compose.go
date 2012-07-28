package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	CompositionPrefix = "[["
	CompositionSuffix = "]]"
)

// ComposeTemplate recursively resolves Compose directives in the given
// template file. A Compose directive is just a naïve "include": it drops the
// content of another template file (specified by filename) in-place into the
// current template.
//
//  a.template: "hello {{name}}"
//  b.template: "<p>[[a.template]]</p>"
//  c.template: "<body>[[b.template]]</body>"
//
// Composing c.template yields "<body><p>hello {{name}}</p></body>"
//
func ComposeTemplate(templatesDir, filename string) ([]byte, error) {
	b, err := ioutil.ReadFile(templatesDir + "/" + filename)
	if err != nil {
		return []byte{}, err
	}

	str := bytes.NewBuffer(b).String()
	out := []byte{}
	pos := 0
	for {

		// Find the next template-inclusion directive
		a := strings.Index(str[pos:], CompositionPrefix)
		if a == -1 {
			out = append(out, []byte(str[pos:])...) // write ending bytes
			break
		}
		out = append(out, []byte(str[pos:pos+a])...) // write in-between bytes
		filenameBegin := pos + a + len(CompositionPrefix)

		// Find its close tag
		b := strings.Index(str[filenameBegin:], CompositionSuffix)
		if b == -1 {
			return nil, fmt.Errorf(
				"unclosed %s at %s:%d",
				CompositionPrefix,
				filename,
				a,
			)
		}
		filenameEnd := filenameBegin + b

		// Compose that template
		includedFile := str[filenameBegin:filenameEnd]
		if strings.IndexAny(includedFile, "\r\n") != -1 {
			return nil, fmt.Errorf(
				"likely unclosed %s at %s:%d",
				CompositionPrefix,
				filename,
				a,
			)
		}
		replacement, err := ComposeTemplate(templatesDir, includedFile)
		if err != nil {
			return nil, err
		}

		// Drop it in-place and repeat
		out = append(out, replacement...)
		pos = filenameEnd + len(CompositionSuffix)
	}

	return out, nil
}
