package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteOutput writes buf to filename.
// It will create any directories required so that filename can exist.
func WriteOutput(buf []byte, filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("%s: short write: %d < %d", filename, n, len(buf))
	}
	return nil
}
