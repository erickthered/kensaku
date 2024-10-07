package main

import (
  "fmt"
//  "flag"
  "os"
  "strings"
  "unicode"
)

var (
  url string
)

type InvertedIndex map[string][]int

func main() {
/*  flag.StringVar(&url, "url", "", "URL to index")
  flag.Parse()
  if (url == "") {
    fmt.Fprintln(os.Stderr, "Flags:")
    flag.PrintDefaults()
    return
  }
  fmt.Println("Hello World!!")
  fmt.Println("Indexing " + url)*/
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

  query := os.Args[1]
  results := SearchIndex(index, query)

  fmt.Println("Search results for: ", query)
  for _, id := range results {
    fmt.Printf("Document %d : %s\n", id, documents[id])
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
