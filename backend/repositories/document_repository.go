package repositories

import (
	"bpt-knowledge-center/backend/config"
	"bpt-knowledge-center/backend/models"
	"fmt"
	"log"
	"os"

	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
)

// SaveDocument persists the parsed document into Couchbase
func SaveDocument(doc *models.Document) error {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")

	collection := config.GetCollection(bucketName, scopeName, "bpt-docs")

	// 2. Generate ID if missing
	if doc.ID == "" {
		doc.ID = "doc::" + uuid.New().String()
	}
	doc.Type = "document"

	// 3. Upsert (Insert or Update)
	_, err := collection.Upsert(doc.ID, doc, &gocb.UpsertOptions{})
	if err != nil {
		log.Printf("Failed to save document to Couchbase: %v", err)
		return err
	}

	return nil
}

func GetAllDocuments() ([]models.Document, error) {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")
	collectionName := os.Getenv("DB_COLLECTION")
	if collectionName == "" {
		collectionName = "bpt-docs"
	}

	// Query just the metadata fields
	query := fmt.Sprintf("SELECT id, filename, uploaded_at, element_count FROM `%s`.`%s`.`%s` WHERE type = 'document' ORDER BY uploaded_at DESC", bucketName, scopeName, collectionName)

	rows, err := config.Cluster.Query(query, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		var doc models.Document
		// FIX: Use .Row() instead of .Scan() for Couchbase
		if err := rows.Row(&doc); err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}
