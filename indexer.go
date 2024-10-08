package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func InsertDocument(db *sql.DB, path string, total_terms int, parser string) int {
	path, _ = filepath.Abs(path)
	res, err := db.Exec("INSERT INTO document(path, total_terms, parser) VALUES (?, ?, ?)", path, total_terms, parser)
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

func IndexPath(sourcePath string) {
	db, err := sql.Open("sqlite3", SEARCH_DB_NAME)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Adding file %s to index...\n", sourcePath)
	content, err := ReadFile(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	// @TODO: obtain the proper parser for sourcePath
	switch path.Ext(strings.ToLower(sourcePath)) {
	case ".html", ".htm":
		IndexHtmlFile(db, sourcePath, content)
	default:
		IndexTextFile(db, sourcePath, content)
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

func IndexTextFile(db *sql.DB, path string, content string) {
	tokens := Tokenize(content)
	unique_tokens := make(map[string][]int)

	docId := InsertDocument(db, path, len(tokens), "text")
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
}

func IndexHtmlFile(db *sql.DB, path string, html_content string) {
	domDoc := html.NewTokenizer(strings.NewReader(html_content))
	previousStartTokenTest := domDoc.Token()
	text_contents := ""

loopDomTest:
	for {
		tt := domDoc.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest // End of the document,  done
		case tt == html.StartTagToken:
			previousStartTokenTest = domDoc.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" || previousStartTokenTest.Data == "style" {
				continue
			}
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDoc.Text())))
			if len(TxtContent) > 0 {
				text_contents += TxtContent + "\n"
			}
		}
	}

	tokens := Tokenize(text_contents)
	unique_tokens := make(map[string][]int)

	docId := InsertDocument(db, path, len(tokens), "html")
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
}
