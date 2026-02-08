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

// ChunkMatch represents a search result with text and source metadata
type ChunkMatch struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Page   int    `json:"page"`
}

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

// 2. The Main Search Function - Returns chunks with source metadata
func SearchSimilarChunks(vectorData []float32) ([]ChunkMatch, error) {
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

	// D. Define Search Options - Include all chunk fields
	opts := &gocb.SearchOptions{
		Limit:  3,
		Fields: []string{"*"}, // Request all fields
	}

	// E. Execute Search (Scoped)
	bucket := config.Cluster.Bucket("bpt-knowledge-center")
	scope := bucket.Scope("knowledge-base")

	result, err := scope.Search("knowledge_vector_search", request, opts)
	if err != nil {
		log.Printf("Search query failed: %v", err)
		return nil, err
	}

	// F. Parse Results with metadata
	var matches []ChunkMatch
	for result.Next() {
		row := result.Row()

		// Get the document ID to fetch full document
		docID := row.ID
		log.Printf("DEBUG: Found document ID: %s", docID)

		var fields map[string]interface{}
		if err := row.Fields(&fields); err == nil {
			log.Printf("DEBUG: Fields returned: %+v", fields)

			match := ChunkMatch{}

			// Try different field patterns
			// Pattern 1: Direct chunks fields
			if val, exists := fields["chunks.text"]; exists {
				match.Text = fmt.Sprintf("%v", val)
			}
			if val, exists := fields["chunks.metadata.source"]; exists {
				match.Source = fmt.Sprintf("%v", val)
			}
			if val, exists := fields["chunks.metadata.page"]; exists {
				if pageNum, ok := val.(float64); ok {
					match.Page = int(pageNum)
				}
			}

			// Pattern 2: Check for nested chunks array
			if chunks, exists := fields["chunks"]; exists {
				log.Printf("DEBUG: chunks field exists: %+v", chunks)
				if chunkArr, ok := chunks.([]interface{}); ok && len(chunkArr) > 0 {
					if chunk, ok := chunkArr[0].(map[string]interface{}); ok {
						if text, exists := chunk["text"]; exists {
							match.Text = fmt.Sprintf("%v", text)
						}
						if meta, exists := chunk["metadata"]; exists {
							if metadata, ok := meta.(map[string]interface{}); ok {
								if source, exists := metadata["source"]; exists {
									match.Source = fmt.Sprintf("%v", source)
								}
								if page, exists := metadata["page"]; exists {
									if pageNum, ok := page.(float64); ok {
										match.Page = int(pageNum)
									}
								}
							}
						}
					}
				}
			}

			// If we still don't have source, try to get from document directly
			if match.Source == "" && docID != "" {
				match = fetchDocumentMetadata(docID, match)
			}

			if match.Text != "" {
				matches = append(matches, match)
			}
		}
	}

	// Check for iterator errors
	if err := result.Err(); err != nil {
		return nil, err
	}

	log.Printf("DEBUG: Total matches with metadata: %+v", matches)
	return matches, nil
}

// Helper to fetch document metadata directly from Couchbase
func fetchDocumentMetadata(docID string, match ChunkMatch) ChunkMatch {
	bucket := config.Cluster.Bucket("bpt-knowledge-center")
	scope := bucket.Scope("knowledge-base")
	collection := scope.Collection("bpt-docs")

	result, err := collection.Get(docID, nil)
	if err != nil {
		log.Printf("DEBUG: Failed to get document %s: %v", docID, err)
		return match
	}

	var doc map[string]interface{}
	if err := result.Content(&doc); err != nil {
		log.Printf("DEBUG: Failed to parse document: %v", err)
		return match
	}

	log.Printf("DEBUG: Full document: %+v", doc)

	// Extract filename from document
	if filename, exists := doc["filename"]; exists {
		match.Source = fmt.Sprintf("%v", filename)
	}

	// Try to get from chunks metadata
	if chunks, exists := doc["chunks"]; exists {
		if chunkArr, ok := chunks.([]interface{}); ok && len(chunkArr) > 0 {
			if chunk, ok := chunkArr[0].(map[string]interface{}); ok {
				if meta, exists := chunk["metadata"]; exists {
					if metadata, ok := meta.(map[string]interface{}); ok {
						if source, exists := metadata["source"]; exists && match.Source == "" {
							match.Source = fmt.Sprintf("%v", source)
						}
						if page, exists := metadata["page"]; exists && match.Page == 0 {
							if pageNum, ok := page.(float64); ok {
								match.Page = int(pageNum)
							}
						}
					}
				}
			}
		}
	}

	return match
}
