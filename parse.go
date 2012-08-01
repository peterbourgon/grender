package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	YYYYMMDD = "([0-9]{4})-([0-9]{2})-([0-9]{2})"
	Title    = "([0-9A-Za-z_-]+)"
)

var (
	R = regexp.MustCompile(fmt.Sprintf("%s-?%s?", YYYYMMDD, Title))
)

// ParseSourceFile reads the given filename (assumed to be a source file, and a
// relative path which must exist  under the passed parentDir) and extracts a
// Context object from its metadata.
//
// If err is nil, the returned Context is guaranteed to contain values for:
//  ckey - content; containing the Markdown-rendered body of the source file
//  tkey - template file that should be used to render the content
//  okey - the output filename that should be rendered-to
//
// If err is nil and the filename matches the blog entry pattern, the returned
// context is guaranteed to contain a value for
//  ikey - the index-tuple entry for this source file
//
func ParseSourceFile(
	parentDir string,
	filename string,
	idx Index,
	bpath string,
	delim string,
	ckey string,
	tkey string,
	ikey string,
	okey string,
	ext string,
) (ctx Context, err error) {
	ctx = make(Context)

	// compose complete filename
	if !strings.HasSuffix(parentDir, "/") {
		parentDir = parentDir + "/"
	}
	absFilename := parentDir + filename

	// read file
	f, err := os.Open(absFilename)
	if err != nil {
		return
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	s := string(buf)

	// separate metadata from content, and dump content to context
	if idx := strings.Index(s, delim); idx >= 0 {
		delimiterCutoff := idx + len(delim) + 1 // plus '\n'
		content := buf[delimiterCutoff:]

		switch strings.ToLower(filepath.Ext(filename)) {
		case ".md":
			content = RenderMarkdown(content)
		}

		ctx[ckey] = strings.TrimSpace(string(content))
		buf = buf[:idx] // buf contains only metadata
	} else {
		ctx[ckey] = "" // no content
	}

	// autopopulate ikey if it looks like a blog entry
	basename := Basename("", filename)
	a := R.FindAllStringSubmatch(basename, -1)
	Logf("%s: %v", basename, a)
	if a != nil && len(a) > 0 && len(a[0]) > 3 {
		Logf("%s: blog entry", basename)
		year, month, day, title := a[0][1], a[0][2], a[0][3], ""
		if len(a[0]) > 3 {
			title = strings.Replace(a[0][4], "-", " ", -1)
			if len(title) > 1 {
				title = strings.ToTitle(title)[:1] + title[1:]
			}
		}
		ctx[ikey] = map[string]string{
			"key":   basename,
			"year":  year,
			"month": month,
			"day":   day,
			"title": title,
			"url":   fmt.Sprintf("%s/%s.%s", bpath, basename, ext),
		}
	} else {
		Logf("%s: not a blog entry", basename)
	}

	// unmarshal metadata as YAML
	if err = goyaml.Unmarshal(buf, ctx); err != nil {
		return
	}

	// check for template key: missing = fatal
	if _, ok := ctx[tkey]; !ok {
		err = fmt.Errorf("%s: '%s' not provided", filename, tkey)
		return
	}

	// check for output file key: missing = need to deduce from basename
	if _, ok := ctx[okey]; !ok {
		ctx[okey] = Basename(parentDir, filename)
	}

	return
}
