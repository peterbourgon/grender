package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"strings"
)

// ParseSourceFile reads the given filename (assumed to be a source file, and a
// relative path which must exist  under the passed parentDir) and extracts a
// Context object from its metadata.
//
// If err is nil, the returned Context is guaranteed to contain entries for ckey
// (the content key, containing the un-rendered body of the source file), tkey
// (the template file that should be used to render the content), and okey (the
// output file that should be rendered).
func ParseSourceFile(
	parentDir string,
	filename string,
	delim string,
	ckey string,
	tkey string,
	okey string,
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
		ctx[ckey] = strings.TrimSpace(string(RenderMarkdown(buf[delimiterCutoff:])))
		buf = buf[:idx] // buf contains only metadata
	} else {
		ctx[ckey] = "" // no content
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
