package main

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

func Debugf(format string, args ...interface{}) {
	if *debug {
		log.Printf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if *verbose {
		log.Printf(format, args...)
	}
}

func Warningf(format string, args ...interface{}) {
	log.Printf("Warning: "+format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf("Fatal: "+format, args...)
}
