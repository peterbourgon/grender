package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"strings"
)

// ParseSourceFile reads the given filename (assumed to be a relative file under
// *sourcePath) and produces a parsed SourceFile object from its contents.
func ParseSourceFile(filename string) (sf *SourceFile, err error) {
	sf = NewSourceFile(filename)

	// read file
	f, err := os.Open(*sourcePath + "/" + filename)
	if err != nil {
		return
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	s := string(buf)

	// separate content
	if idx := strings.Index(s, *metadataDelimiter); idx >= 0 {
		delimiterCutoff := idx + len(*metadataDelimiter) + 1 // plus '\n'
		contentBuf := buf[delimiterCutoff:]

		switch strings.ToLower(filepath.Ext(filename)) {
		case ".md":
			contentBuf = RenderMarkdown(contentBuf)
		}

		sf.Metadata[*contentKey] = strings.TrimSpace(string(contentBuf))
		buf = buf[:idx] // buf shall contain only metadata
	}

	// if the filename looks like a blog entry, autopopulate some metadata
	if y, m, d, t, err := sf.BlogEntry(); err == nil {
		sf.Metadata["output"] = fmt.Sprintf("%s/%s", *blogPath, sf.Basename)
		sf.Metadata["sortkey"] = sf.Basename
		sf.Metadata["year"] = y
		sf.Metadata["month"] = m
		sf.Metadata["day"] = d
		sf.Metadata["title"] = t
		sf.Metadata["url"] = fmt.Sprintf("%s.%s", sf.Output(), *outputExtension)
	}

	// read remaining metadata as YAML
	if err = goyaml.Unmarshal(buf, sf.Metadata); err != nil {
		return
	}

	// check for template key: missing = fatal
	if sf.Template() == "" {
		err = fmt.Errorf("%s: '%s' not provided", filename, *templateKey)
		return
	}

	// check for output file key: missing = need to deduce from basename
	if sf.Output() == "" {
		sf.Metadata[*outputKey] = Basename(*sourcePath, filename)
	}

	err = nil // just in case
	return
}
