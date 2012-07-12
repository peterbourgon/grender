package main

import (
	"io/ioutil"
	"bytes"
	"strings"
	"fmt"
	. "github.com/peterbourgon/bonus/xlog"
)

const (
	CompositionPrefix = "[["
	CompositionSuffix = "]]"
)

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
		Debugf("ComposeTemplate %s: pos=%d", filename, pos)
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
		Debugf("%s: %s @ %d, %s @ %d", filename, CompositionPrefix, filenameBegin, CompositionSuffix, filenameEnd)

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
		Debugf("%s: includedFile @ %d: %s", filename, filenameBegin, includedFile)
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
