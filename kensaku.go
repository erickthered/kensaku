package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type InvertedIndex map[string][]int

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: kensaku <serve|index>")
		return
	}

	switch os.Args[1] {
	case "serve":
		documents := []string{
			"Go is an open source programming language.",
			"Go makes it easy to build simple, reliable, and efficient software.",
			"The Go programming language is efficient and fast.",
		}

		index := make(InvertedIndex)

		for i, document := range documents {
			tokens := Tokenize(document)
			AddToIndex(index, i, tokens)
		}

		// Add HTTP library to handle search on indexed documents
		fmt.Println("Initializing web server on port 8080")
		http.HandleFunc("/search", searchHandler(index, documents))
		http.ListenAndServe(":8080", nil)
	default:
		fmt.Println("Usage: kensaku <serve|index>")
	}
}

func searchHandler(index InvertedIndex, documents []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		results := SearchIndex(index, query)

		fmt.Fprintf(w, "Search results for: %s\n", query)
		for _, id := range results {
			fmt.Fprintf(w, "Document %d: %s\n", id, documents[id])
		}
	}
}

// function to tokenize a string
func Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
}

// function to add to the reverse index
func AddToIndex(index InvertedIndex, docID int, tokens []string) {
	for _, token := range tokens {
		token = strings.ToLower(token)
		index[token] = append(index[token], docID)
	}
}

// function to search the index for a given word
func SearchIndex(index InvertedIndex, query string) []int {
	tokens := Tokenize(query)
	docIds := make(map[int]int)

	for _, token := range tokens {
		if ids, found := index[strings.ToLower(token)]; found {
			for _, id := range ids {
				docIds[id]++
			}
		}
	}

	var result []int
	for id, count := range docIds {
		if count == len(tokens) {
			result = append(result, id)
		}
	}

	return result
}
