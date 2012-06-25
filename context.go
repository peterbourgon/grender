package main

import (
	// "github.com/kylelemons/go-gypsy/yaml"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Context map[string]interface{}

const (
// TODO
)

func ContextFrom(sourceFile string) (string, Context, string, error) {

	templateFile, ctx, outputFile := "", Context{}, ""

	f, err := os.Open(sourceFile)
	if err != nil {
		return templateFile, ctx, outputFile, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {

		line, isPrefix, err := r.ReadLine()
		if err != nil {
			return templateFile, ctx, outputFile, err
		}
		if isPrefix {
			return templateFile, ctx, outputFile, fmt.Errorf("%s: too big", sourceFile)
		}
		s := strings.TrimSpace(string(line))
		fmt.Printf("%s", s) // TODO
	}

	return templateFile, ctx, outputFile, nil
}
