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
	stack := map[string]map[string]interface{}{}
	return func(path string, info os.FileInfo, _ error) error {
		fmt.Printf("entering '%s', stack=%v\n", path, stack)
		if info.IsDir() {
			return nil // descend
		}

		switch filepath.Ext(path) {
		case ".json":
			m, err := readJSON(path)
			if err != nil {
				fmt.Printf("%s: %s\n", path, err)
				os.Exit(2)
			}
			stack[filepath.Dir(path)] = m

		case ".htm", ".html":
			copyFile(filepath.Join(*targetDir, diffPath(*sourceDir, path)), path)

		default:
			fmt.Printf("%s: no action\n", path)
		}

		return nil
	}
}
