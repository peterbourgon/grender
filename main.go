package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	sourceDir = flag.String("source", "src", "path to site source (input)")
	targetDir = flag.String("target", "tgt", "path to site target (output)")
)

func init() {
	flag.Parse()

	var err error
	for _, s := range []*string{sourceDir, targetDir} {
		if *s, err = filepath.Abs(*s); err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
	}
}

func main() {
	filepath.Walk(*sourceDir, sourceWalk())
}

func sourceWalk() filepath.WalkFunc {
	s := NewStack()

	readAndAdd := func(path string) {
		m, err := readJSON(path)
		if err != nil {
			fmt.Printf("%s: %s\n", path, err)
			os.Exit(2)
		}
		s.Add(filepath.Dir(path), m)
		fmt.Printf("%-70s added to stack: %v\n", path, m)
	}

	cp := func(path string) {
		dst := filepath.Join(*targetDir, diffPath(*sourceDir, path))
		copyFile(dst, path)
		fmt.Printf("%-70s copied to %s\n", path, dst)
	}

	ext := map[string]func(path string){
		".json": readAndAdd,
		".htm":  cp,
		".html": cp,
	}

	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // descend
		}

		if f, ok := ext[filepath.Ext(path)]; ok {
			f(path)
		} else {
			fmt.Printf("%-70s no action\n", path)
		}
		return nil
	}
}
