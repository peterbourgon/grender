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
		outputPath := fmt.Sprintf(
			"%s/%4d/%02d/%02d/%s.%s",
			*blogPath,
			y,
			m,
			d,
			Display2Filename(t),
			*outputExtension,
		)

		sf.Metadata[*outputKey] = outputPath
		sf.Metadata[*sortkeyKey] = sf.Basename
		sf.Metadata[YearKey] = fmt.Sprintf("%4d", y)
		sf.Metadata[MonthKey] = fmt.Sprintf("%02d", m)
		sf.Metadata[DayKey] = fmt.Sprintf("%02d", d)
		sf.Metadata[TitleKey] = t
		sf.Metadata[URLKey] = outputPath
	}

	// read remaining metadata as YAML
	if err = goyaml.Unmarshal(buf, sf.Metadata); err != nil {
		return
	}

	// check for some keys
	if sf.getString(*templateKey) == "" {
		err = fmt.Errorf("%s: '%s' not provided", filename, *templateKey)
		return
	}
	if sf.getString(*outputKey) == "" {
		sf.Metadata[*outputKey] = fmt.Sprintf(
			"%s.%s",
			Basename(*sourcePath, filename),
			*outputExtension,
		)
	}

	err = nil // just in case
	return
}
