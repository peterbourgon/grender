package main

import (
	"log"
)

func Debugf(format string, args ...interface{}) {
	if *debug {
		log.Printf("DEBUG "+format, args...)
	}
}

func Problemf(format string, args ...interface{}) {
	log.Printf("PROBLEM "+format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf("FATAL "+format, args...)
}
