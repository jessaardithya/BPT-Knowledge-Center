package controllers

import (
	"bpt-knowledge-center/backend/models"
	"bpt-knowledge-center/backend/repositories"
	"bpt-knowledge-center/backend/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadDocument(c *gin.Context) {
	// 1. Get the file from the request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Get optional display name
	displayName := c.PostForm("display_name")
	if displayName == "" {
		displayName = fileHeader.Filename
	}

	// Check if this is a re-upload (update existing document)
	documentID := c.PostForm("document_id")
	var existingVersion int = 0

	// If re-uploading, get existing version and DELETE old document first
	// This removes old chunks to prevent AI conflicts
	if documentID != "" {
		existingDoc, err := repositories.GetDocumentByID(documentID)
		if err == nil && existingDoc != nil {
			existingVersion = existingDoc.Version
			// Delete old document to remove old chunks from vector index
			if err := repositories.DeleteDocument(documentID); err != nil {
				log.Printf("Warning: Failed to delete old document: %v", err)
			}
		}
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
		return
	}
	defer file.Close()

	fileURL, err := services.UploadFile(file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage upload failed: " + err.Error()})
		return
	}

	tempPath := "./temp/" + fileHeader.Filename
	if err := c.SaveUploadedFile(fileHeader, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temp file"})
		return
	}

	parsedData, err := services.SendToParser(tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Parsing failed: " + err.Error()})
		return
	}

	// Create new document (with same ID if re-upload, new ID otherwise)
	var docID string
	var version int = 1

	if documentID != "" && existingVersion > 0 {
		docID = documentID
		version = existingVersion + 1
	} else {
		docID = "doc::" + uuid.New().String()
	}

	doc := models.Document{
		ID:           docID,
		Type:         "document",
		Filename:     fileHeader.Filename,
		DisplayName:  displayName,
		FileURL:      fileURL,
		ContentType:  parsedData.ContentType,
		UploadedAt:   time.Now(),
		UpdatedAt:    time.Now(),
		ElementCount: parsedData.ElementCount,
		Version:      version,
		DocType:      "knowledge-base.bpt-docs",
		Chunks:       make([]models.DocumentChunk, len(parsedData.Data)),
	}

	for i, item := range parsedData.Data {
		doc.Chunks[i] = models.DocumentChunk{
			ChunkID:  item.ElementID,
			Text:     item.Text,
			Type:     item.Type,
			Metadata: item.Metadata,
			Vector:   item.Vector,
		}
	}

	// Save new version to Couchbase
	if err := repositories.SaveDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database save failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ingestion successful",
		"id":      doc.ID,
		"url":     doc.FileURL,
		"count":   doc.ElementCount,
		"version": doc.Version,
	})
}
