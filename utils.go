package main

import (
	"regexp"
)

// function to tokenize a string
func Tokenize(text string) []string {
	re := regexp.MustCompile(`\w+('[d,m,s,t])?`)
	return re.FindAllString(text, -1)
}
