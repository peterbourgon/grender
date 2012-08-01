package main

type IndexTuple map[string]string

type Index map[string][]IndexTuple

func NewIndex() Index {
	return Index{}
}
