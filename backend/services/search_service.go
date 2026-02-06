package services

import (
	"bpt-knowledge-center/backend/config"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/couchbase/gocb/v2"
	"github.com/couchbase/gocb/v2/vector"
)

// 1. Helper to call Python and get the Vector
func GetQueryVector(question string) ([]float32, error) {
	// Payload
	requestBody, _ := json.Marshal(map[string]string{
		"text": question,
	})

	// Call Python Service
	resp, err := http.Post("http://localhost:8000/api/v1/embed", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Vector []float32 `json:"vector"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Vector, nil
}

// 2. The Main Search Function
func SearchSimilarChunks(vectorData []float32) ([]string, error) {
	// A. Define Vector Query
	// Matches the "vector" field inside the "chunks" nested array
	vQuery := vector.NewQuery("chunks.vector", vectorData).
		NumCandidates(3)

	// B. Define Vector Search
	vSearch := vector.NewSearch([]*vector.Query{vQuery}, nil)

	// C. Define Search Request
	request := gocb.SearchRequest{
		VectorSearch: vSearch,
	}

	// D. Define Search Options
	opts := &gocb.SearchOptions{
		Limit:  3,
		Fields: []string{"chunks.text"},
	}

	// E. Execute Search (Scoped)
	// The index is scoped to 'bpt-knowledge-center' -> 'knowledge-base'
	bucket := config.Cluster.Bucket("bpt-knowledge-center")
	scope := bucket.Scope("knowledge-base")

	result, err := scope.Search("knowledge_vector_search", request, opts)
	if err != nil {
		log.Printf("Search query failed: %v", err)
		return nil, err
	}

	// D. Parse Results
	var matches []string
	for result.Next() {
		row := result.Row()

		var fields map[string]interface{}
		if err := row.Fields(&fields); err == nil {
			if val, exists := fields["chunks.text"]; exists {
				matches = append(matches, fmt.Sprintf("%v", val))
			}
		}
	}

	// Check for iterator errors
	if err := result.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}
