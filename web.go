package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// SearchHandler for web requests
func SearchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		start := time.Now()
		results := SearchIndex(db, query)
		duration := time.Since(start)

		var documentId int
		var documentPath string
		var documentRanking float32

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{ \"result_count\": 10, \"results\": [")
		isFirst := true
		for results.Next() {
			results.Scan(&documentId, &documentPath, &documentRanking)
			if !isFirst {
				fmt.Fprintf(w, ",")
			}
			fmt.Fprintf(w, `{ "document_id" : %d, "document": "%s", "ranking": %f }`, documentId, documentPath, documentRanking*100)
			isFirst = false
		}
		fmt.Fprintf(w, "], \"duration\": %f }", duration.Seconds())
	}
}

func DocumentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docId := r.PathValue("doc_id")

		var documentPath string
		var documentParser string
		var documentTerms int

		err := db.QueryRow("SELECT path, parser, total_terms FROM document WHERE id=?", docId).Scan(&documentPath, &documentParser, &documentTerms)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Fprintf(w, "Document %s not found.\n", docId)
				return
			} else {
				fmt.Fprintf(w, "Error: %s", err)
				return
			}
		}

		res, err := db.Query("SELECT t.term, dt.term_count, 1.0*dt.term_count/d.total_terms FROM docterm dt, term t, document d WHERE dt.document_id = d.id AND dt.term_id = t.id AND dt.document_id = ? ORDER BY dt.term_count DESC, t.term ASC", docId)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id": %s, "path": "%s", "parser": "%s", "total_terms": %d, "keywords": [`, docId, documentPath, documentParser, documentTerms)

		var keyword string
		var keyword_count int
		var keyword_density float32
		var isFirst = true
		for res.Next() {
			res.Scan(&keyword, &keyword_count, &keyword_density)
			if !isFirst {
				fmt.Fprint(w, ",\n")
			}
			fmt.Fprintf(w, `{"keyword": "%s", "keyword_count": %d, "keyword_density": %f}`, keyword, keyword_count, keyword_density)
			isFirst = false
		}
		fmt.Fprint(w, `]}`)
	}
}
