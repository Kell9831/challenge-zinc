package handlers

import (
	"io"
	"net/http"
)

type Email struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Body    string `json:"body"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
		return
	}

	page, size, err := getPaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from := (page - 1) * size

	zincQuery, err := createZincQuery(body, from, size)
	if err != nil {
		http.Error(w, "error creating ZincSearch query", http.StatusInternalServerError)
		return
	}

	zincResponse, err := queryZincSearch(zincQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, totalResults, err := processZincResponse(zincResponse)
	if err != nil {
		http.Error(w, "Error processing ZincSearch response", http.StatusInternalServerError)
		return
	}

	formatAndReturnResponse(w, results,totalResults, page, size)
}

