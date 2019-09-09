package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var replacer *strings.Replacer

func init() {
	replacer = strings.NewReplacer("\r", "", "\n", "", "\t", "")
}

func strip(val string) string {
	return replacer.Replace(strings.TrimSpace(val))
}

// Map is map
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// Sometimes is very usefull to have ability dump response to the file
func dumpResponse(filename string, val io.ReadCloser) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Can't create file %s, error: %s\n", filename, err.Error())
	}
	defer file.Close()

	written, err := io.Copy(file, val)
	if err != nil {
		log.Fatalf("Can't save file %s, error: %s\n", filename, err.Error())
	}

	fmt.Printf("\tfile %s (%d bytes) created\n",
		filename,
		written)
}
