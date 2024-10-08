package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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
		IndexPath(os.Args[2])

	case "serve":
		db, err := sql.Open("sqlite3", SEARCH_DB_NAME)
		if err != nil {
			log.Fatal(err)
		}

		fs := http.FileServer(http.Dir("./web"))

		mux := http.NewServeMux()
		mux.HandleFunc("/search", SearchHandler(db))
		mux.HandleFunc("/document/{doc_id}", DocumentHandler(db))
		mux.Handle("/", fs)

		// Add HTTP library to handle search on indexed documents
		fmt.Println("Initializing web server on port 8080...")
		http.ListenAndServe(":8080", mux)

	default:
		fmt.Println("Usage: kensaku <serve|index>")
	}
}

// function to search the index for a given word
func SearchIndex(db *sql.DB, query string) *sql.Rows {
	fmt.Println("query:", query)
	tokens := Tokenize(query)
	query = "SELECT d.id as document_id, d.path, sum(1.0*dt.term_count/d.total_terms) AS ranking FROM "
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
