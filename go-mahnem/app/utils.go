package main

import (
	"log"
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
