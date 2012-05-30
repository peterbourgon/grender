package main

import (
	"testing"
	"os"
	"path"
	"io/ioutil"
)

func writeTestHierarchy(t *testing.T) string {
	root, err := ioutil.TempDir(os.TempDir(), "grender-walk-test")
	if err != nil {
		t.Fatalf("creating temp dir: %s", err)
	}

	// root/                    should become
	//  |2012/                  |2012/
	//  |  |2.txt               |   |3/
	//  |  |3/                  |   |   |deep.txt <---------1
	//  |     |deep.txt         |   |2.txt <----------------2
	//  |_.txt                  |subdir/
	//  |foo.txt                |   |subsubdir/
	//  |subdir/                |   |   |deepdir/
	//     |subdir_file.txt     |   |       |deepest.txt <--3
	//     |subsubdir/          |   |subdir_file.txt <------4
	//        |deepdir/         |_.txt <--------------------5 <-lexigraphical
	//           |deepest.txt   |foo.txt <------------------6 <-ordering

	ioutil.WriteFile(root+"/_.txt", []byte{}, 0755)
	ioutil.WriteFile(root+"/foo.txt", []byte{}, 0755)
	os.MkdirAll(root+"/2012", 0755)
	ioutil.WriteFile(root+"/2012/2.txt", []byte{}, 0755)
	os.MkdirAll(root+"/2012/3", 0755)
	ioutil.WriteFile(root+"/2012/3/deep.txt", []byte{}, 0755)
	os.MkdirAll(root+"/subdir/subsubdir/deepdir", 0755)
	ioutil.WriteFile(root+"/subdir/subdir_file.txt", []byte{}, 0755)
	ioutil.WriteFile(root+"/subdir/subsubdir/deepdir/deepest.txt", []byte{}, 0755)

	return root
}

func TestWalkAllFiles(t *testing.T) {
	root := writeTestHierarchy(t)
	defer os.RemoveAll(root)

	files := []string{}
	visitor := func(name string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path.Base(name))
		}
		return nil // descend dirs
	}
	if err := DepthFirstWalk(root, visitor); err != nil {
		t.Errorf("visitor gave error: %s", err)
	}

	expectedFiles := []string{"deep.txt", "2.txt", "deepest.txt", "subdir_file.txt", "_.txt", "foo.txt"}
	if len(files) != len(expectedFiles) {
		t.Fatalf("expected %d files, got %d: %v", len(expectedFiles), len(files), files)
	}
	for i := 0; i < len(files); i++ {
		if expectedFiles[i] != files[i] {
			t.Errorf("visit %d: expected '%s', got '%s'", i, expectedFiles[i], files[i])
		}
	}
}

func TestChoosyWalker(t *testing.T) {
	root := writeTestHierarchy(t)
	defer os.RemoveAll(root)

	chosen := ""
	chooser := func(name string, info os.FileInfo, err error) error {
		if chosen == "" && !info.IsDir() {
			chosen = path.Base(name) // first one is chosen
		}
		return nil
	}
	if err := DepthFirstWalk(root, chooser); err != nil {
		t.Errorf("chooser gave error: %s", err)
	}
	expected := "deep.txt"
	if chosen != expected {
		t.Errorf("expected '%s', got '%s'", expected, chosen)
	}
}
