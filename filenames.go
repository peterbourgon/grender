package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Filenames walks the path rooted at 'root' and returns a slice containing
// the (relative) filenames of every file under that path.
func Filenames(root string) []string {
	root = filepath.Clean(root)
	rootCutoff := len(root)
	if root[rootCutoff-1:] != "/" {
		rootCutoff++ // also cut the trailing '/' in any discovered paths
	}

	filenames := []string{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Logf("Filenames under %s: %s", root, err)
			return err // abort
		}
		if info.IsDir() {
			return nil // descend
		}
		if !strings.HasPrefix(path, root) {
			return fmt.Errorf("source file '%s': missing '%s'", path, root)
		}
		path = path[rootCutoff:]
		filenames = append(filenames, path)
		return nil
	})
	return filenames
}

// Basename returns filename stripped of its extension.
// If the filename begins with prefix, that is also stripped.
func Basename(prefix, filename string) string {
	prefix, filename = path.Clean(prefix), path.Clean(filename)
	endCutoff := len(filename) - len(path.Ext(filename))
	startCutoff := 0
	if strings.HasPrefix(filename, prefix) {
		startCutoff = len(prefix) + 1
	}
	return filename[startCutoff:endCutoff]
}
