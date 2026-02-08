package repositories

import (
	"bpt-knowledge-center/backend/config"
	"bpt-knowledge-center/backend/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
)

// SaveDocument persists the parsed document into Couchbase
func SaveDocument(doc *models.Document) error {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")

	collection := config.GetCollection(bucketName, scopeName, "bpt-docs")

	// Generate ID if missing
	if doc.ID == "" {
		doc.ID = "doc::" + uuid.New().String()
	}
	doc.Type = "document"

	// Set version to 1 if new document
	if doc.Version == 0 {
		doc.Version = 1
	}

	// Set display name to filename if not provided
	if doc.DisplayName == "" {
		doc.DisplayName = doc.Filename
	}

	// Upsert (Insert or Update)
	_, err := collection.Upsert(doc.ID, doc, &gocb.UpsertOptions{})
	if err != nil {
		log.Printf("Failed to save document to Couchbase: %v", err)
		return err
	}

	return nil
}

func UpdateDocumentMetadata(id string, category string, description string) error {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")
	collectionName := os.Getenv("DB_COLLECTION")

	query := fmt.Sprintf("UPDATE `%s`.`%s`.`%s` SET category = $1, description = $2 WHERE id = $3", bucketName, scopeName, collectionName)

	_, err := config.Cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{category, description, id},
	})

	return err
}

// UpdateDocumentName updates the display name of a document
func UpdateDocumentName(id string, displayName string) error {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")
	collectionName := os.Getenv("DB_COLLECTION")
	if collectionName == "" {
		collectionName = "bpt-docs"
	}

	query := fmt.Sprintf("UPDATE `%s`.`%s`.`%s` SET display_name = $1 WHERE id = $2", bucketName, scopeName, collectionName)

	_, err := config.Cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{displayName, id},
	})

	return err
}

// GetDocumentByID retrieves a single document by ID
func GetDocumentByID(id string) (*models.Document, error) {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")
	collection := config.GetCollection(bucketName, scopeName, "bpt-docs")

	result, err := collection.Get(id, nil)
	if err != nil {
		return nil, err
	}

	var doc models.Document
	if err := result.Content(&doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

// IncrementDocumentVersion increments version and updates the document
func IncrementDocumentVersion(id string) (int, error) {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")
	collectionName := os.Getenv("DB_COLLECTION")
	if collectionName == "" {
		collectionName = "bpt-docs"
	}

	// Get current version
	doc, err := GetDocumentByID(id)
	if err != nil {
		return 0, err
	}

	newVersion := doc.Version + 1
	now := time.Now()

	query := fmt.Sprintf("UPDATE `%s`.`%s`.`%s` SET version = $1, updated_at = $2 WHERE id = $3", bucketName, scopeName, collectionName)

	_, err = config.Cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{newVersion, now, id},
	})

	return newVersion, err
}

func DeleteDocument(id string) error {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")
	collection := config.GetCollection(bucketName, scopeName, "bpt-docs")

	_, err := collection.Remove(id, nil)
	return err
}

func GetAllDocuments() ([]models.Document, error) {
	bucketName := os.Getenv("DB_BUCKET")
	scopeName := os.Getenv("DB_SCOPE")
	collectionName := os.Getenv("DB_COLLECTION")
	if collectionName == "" {
		collectionName = "bpt-docs"
	}

	// Query with new fields
	query := fmt.Sprintf("SELECT id, filename, display_name, uploaded_at, updated_at, element_count, version, category, description FROM `%s`.`%s`.`%s` WHERE type = 'document' ORDER BY uploaded_at DESC", bucketName, scopeName, collectionName)
	rows, err := config.Cluster.Query(query, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		var doc models.Document
		if err := rows.Row(&doc); err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}
