package main

import (
	"fmt"
)

type Context map[string]interface{}

func (ctx Context) GetString(key string) (string, error) {
	i, ok := ctx[key]
	if !ok {
		return "", fmt.Errorf("%s: not found", key)
	}
	s, ok := i.(string)
	if !ok {
		return "", fmt.Errorf("%s: not a string", key)
	}
	return s, nil
}
