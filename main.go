package main

import (
	"fmt"
	"net/http"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	// Assume we have a function getResults that returns the results and the limit
	results, limit := getResults(r)
	resultCount := len(results)

	// Check if results are truncated
	if resultCount > limit {
		http.Error(w, "Results truncated", http.StatusBadRequest)
		return
	}

	// If resultCount is equal to limit, do not trigger truncation warning
	// Respond with the results
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Results: %v", results)
}

func main() {
	http.HandleFunc("/api/v1/query", queryHandler)
	fmt.Println("Hello, Bounty Hunter!")
	http.ListenAndServe(":8080", nil)
}