package main

import (
	"strings"
	"unicode"
)

// function to tokenize a string
func Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
}
