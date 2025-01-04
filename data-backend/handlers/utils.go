package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getPaginationParams(r *http.Request) (int, int, error) {
	query := r.URL.Query()
	pageParam := query.Get("page")
	sizeParam := query.Get("size")
	page := 1
	size := 10

	if pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil || page <= 0 {
			return 0, 0, fmt.Errorf("invalid 'page' parameter")
		}
	}

	if sizeParam != "" {
		var err error
		size, err = strconv.Atoi(sizeParam)
		if err != nil || size <= 0 {
			return 0, 0, fmt.Errorf("invalid 'size' parameter")
		}
	}

	return page, size, nil
}

func createZincQuery(body []byte, from, size int) (string, error) {
	return fmt.Sprintf(`{
		"search_type": "match",
		"query": %s,
		"from": %d,
		"max_results": %d,
		"_source": ["subject", "from", "to", "body"],
		"sort_fields": ["-@timestamp"]
	}`, string(body), from, size), nil
}

func queryZincSearch(query string) ([]byte, error) {
	zincURL := os.Getenv("ZINC_URL")
	if zincURL == "" {
		return nil, fmt.Errorf("ZINC_URL environment variable is not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", zincURL, strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("error creating request to ZincSearch: %v", err)
	}

	username := os.Getenv("ZINC_USER")
	password := os.Getenv("ZINC_PASSWORD")

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error connecting to ZincSearch: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ZincSearch error: %v", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func processZincResponse(response []byte) ([]Email, int,error) {
	var results struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source Email `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err := json.Unmarshal(response, &results)
	if err != nil {
		return nil,0, fmt.Errorf("error processing response from ZincSearch: %v", err)
	}

	var formattedResults []Email
	for _, hit := range results.Hits.Hits {
		formattedResults = append(formattedResults, hit.Source)
	}

	return formattedResults,results.Hits.Total.Value, nil
}

func formatAndReturnResponse(w http.ResponseWriter, results []Email,totalResults, page, size int) {
	totalPages := (totalResults + size - 1) / size

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, _ := json.Marshal(struct {
		Results     []Email `json:"results"`
		Total       int     `json:"total_results"`
		Page        int     `json:"page"`
		TotalPages  int     `json:"total_pages"`
		ResultsSize int     `json:"results_per_page"`
	}{
		Results: results, 
		Total: totalResults, 
		Page: page, 
		TotalPages: totalPages, 
		ResultsSize: size,
	})

	w.Write(jsonResponse)
}
