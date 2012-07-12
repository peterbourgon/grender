package main

import (
	"fmt"
	"path/filepath"
	"os"
	"regexp"
	"strings"
	"sort"
)

var (
	PostFilenameRegexp *regexp.Regexp
)

const (
	// PostFilenameRE = `([0-9]{4})-([0-9]{2})-([0-9]{2})[-]?([0-9A-Za-z\-]*)`
	PostFilenameRE = `^([0-9]{4})-([0-9]{2})-([0-9]{2})$`
)

func init() {
	PostFilenameRegexp = regexp.MustCompile(PostFilenameRE)
}

// FilesIn returns a slice of files in the given directory.
func FilesIn(dir string) []string {
	paths := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Problemf("%s", err)
			return nil // continue
		}
		if info.IsDir() {
			return nil // descend
		}
		paths = append(paths, path)
		return nil
	})
	Debugf("FilesIn(%s) %v", dir, paths)
	return paths
}

//
//
//

type Post struct {
	Filename string
	Year     string
	Month    string
	Day      string
	Title    string
}

func (p *Post) SortKey() string {
	return fmt.Sprintf("%s%s%s", p.Year, p.Month, p.Day)
}

type Posts []Post

func (a Posts) Len() int           { return len(a) }
func (a Posts) Less(i, j int) bool { return a[i].SortKey() > a[j].SortKey() }
func (a Posts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// PostsIn returns a slice of posts in the given directory.
// They will be in date order, newest first.
func PostsIn(dir string, urlBase string) Posts {
	posts := Posts{}
	for _, path := range FilesIn(dir) {
		if !PostFilenameRegexp.MatchString(filepath.Base(path)) {
			Problemf("%s: invalid filename structure", path)
			continue
		}
		toks := PostFilenameRegexp.FindStringSubmatch(filepath.Base(path))
		if len(toks) != 1+3 {
			Problemf("%s: bad regex parse (%d)", path, len(toks))
			continue
		}
		// Debugf("%s: year=%s month=%s day=%s title='%s'", path, y, m, d, t)
		y, m, d := toks[1], toks[2], toks[3]
		posts = append(posts, Post{path, y, m, d, ""})
	}
	sort.Sort(posts)
	return posts
}

func Titleize(s string) string {
	if s == "" {
		return s
	}
	s = strings.Replace(s, "-", " ", -1)
	return string(strings.ToTitle(s)[0]) + s[1:]
}
