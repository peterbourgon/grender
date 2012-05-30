package main

import (
	"path/filepath"
	"os"
)

func DepthFirstWalk(root string, f filepath.WalkFunc) error {
	type Entry struct {
		path string
		info os.FileInfo
	}

	// Collect & sort entries for this directory
	files, dirs := []Entry{}, []Entry{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path == root {
			return nil
		}
		e := Entry{path, info}
		if info.IsDir() {
			dirs = append(dirs, e)
			return filepath.SkipDir // don't descend right now
		}
		files = append(files, e)
		return nil
	})

	// Process them with the passed WalkFunc, directories first
	for _, entry := range dirs {
		err := f(entry.path, entry.info, nil)
		if err == nil {
			if entry.path == root {
				continue // OK
			}
			err = DepthFirstWalk(entry.path, f) // descend (recurse)
			if err != nil {
				return err
			}
		} else if err == filepath.SkipDir {
			continue // don't descend
		} else {
			return err // abort
		}
	}
	for _, entry := range files {
		if err := f(entry.path, entry.info, nil); err != nil {
			return err // abort
		}
	}
	return nil
}
