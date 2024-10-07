package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

type InvertedIndex map[string][]int

const SEARCH_DB_NAME string = "db/kensaku.db"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: kensaku <serve|index>")
		return
	}

	switch os.Args[1] {
	case "index":
		if len(os.Args) < 3 {
			fmt.Println("Usage: kensaku index <path>")
			return
		}

		index := make(InvertedIndex)
		IndexPath(index, os.Args[2])

	case "serve":
		db, err := sql.Open("sqlite3", SEARCH_DB_NAME)
		if err != nil {
			log.Fatal(err)
		}

		fs := http.FileServer(http.Dir("./web"))

		mux := http.NewServeMux()
		mux.HandleFunc("/search", SearchHandler(db))
		mux.Handle("/", fs)

		// Add HTTP library to handle search on indexed documents
		fmt.Println("Initializing web server on port 8080...")
		http.ListenAndServe(":8080", mux)

	default:
		fmt.Println("Usage: kensaku <serve|index>")
	}
}

func InsertDocument(db *sql.DB, path string, total_terms int) int {
	path, _ = filepath.Abs(path)
	res, err := db.Exec("INSERT INTO document(path, total_terms) VALUES (?, ?)", path, total_terms)
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return int(id)
}

func InsertTerm(db *sql.DB, term string) int {
	term = strings.ToLower(term)
	var existing_id int
	err := db.QueryRow("SELECT id FROM term WHERE term = ?", term).Scan(&existing_id)
	if err == nil {
		return existing_id
	}
	res, err := db.Exec("INSERT INTO term(term) VALUES (?)", term)
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return int(id)
}

func InsertDocterm(db *sql.DB, document_id int, term_id int, positions []int) {
	res, err := db.Exec("INSERT INTO docterm(document_id, term_id, term_count) VALUES (?, ?, ?)", document_id, term_id, len(positions))
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	for _, position := range positions {
		_, err := db.Exec("INSERT INTO docterm_position(docterm_id, position) VALUES (?, ?)", id, position)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func IndexPath(index InvertedIndex, path string) {
	db, err := sql.Open("sqlite3", SEARCH_DB_NAME)
	if err != nil {
		log.Fatal(err)
	}

	// @TODO: obtain the proper parser for path
	fmt.Printf("Adding file %s to index...\n", path)
	content, err := ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	tokens := Tokenize(content)
	unique_tokens := make(map[string][]int)

	docId := InsertDocument(db, path, len(tokens))
	for i, token := range tokens {
		val, err := unique_tokens[token]
		if err {
			unique_tokens[token] = append(val, i)
		} else {
			unique_tokens[token] = []int{i}
		}
	}

	for unique_token, positions := range unique_tokens {
		termId := InsertTerm(db, unique_token)
		InsertDocterm(db, docId, termId, positions)
	}

	db.Close()
}

// Function to read a file and return its contents as a string
func ReadFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// SearchHandler for web requests
func SearchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		start := time.Now()
		results := SearchIndex(db, query)
		duration := time.Since(start)

		var documentPath string
		var documentRanking float32

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{ \"result_count\": 10, \"results\": [")
		isFirst := true
		for results.Next() {
			results.Scan(&documentPath, &documentRanking)
			if !isFirst {
				fmt.Fprintf(w, ",")
			}
			fmt.Fprintf(w, "{ \"document\": \"%s\", \"ranking\": %f }", documentPath, documentRanking*100)
			isFirst = false
		}
		fmt.Fprintf(w, "], \"duration\": %f }", duration.Seconds())
	}
}

// function to tokenize a string
func Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
}

// function to search the index for a given word
func SearchIndex(db *sql.DB, query string) *sql.Rows {
	fmt.Println("query:", query)
	tokens := Tokenize(query)
	query = "SELECT d.path, sum(1.0*dt.term_count/d.total_terms) AS ranking FROM "
	query += "docterm dt, document d, term t WHERE "
	query += "dt.document_id = d.id AND dt.term_id = t.id AND ("

	for i, token := range tokens {
		query += "t.term = '" + strings.ToLower(token) + "'"
		if i < len(tokens)-1 {
			query += " OR "
		}
	}
	query += ") GROUP BY d.path ORDER BY ranking DESC"

	res, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
