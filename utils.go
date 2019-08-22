package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func strip(val string) string {
	re, err := regexp.Compile(`\r?\n`)
	if err != nil {
		log.Fatal(err)
	}
	val = strings.TrimSpace(strings.ReplaceAll(val, "\t", ""))
	return re.ReplaceAllString(val, "")
}

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
