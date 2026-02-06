package controllers

import (
	"bpt-knowledge-center/backend/models"
	"bpt-knowledge-center/backend/repositories"
	"bpt-knowledge-center/backend/services"
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

	// 4. Create Document Model
	doc := models.Document{
		ID:           "doc::" + uuid.New().String(),
		Type:         "document",
		Filename:     fileHeader.Filename,
		FileURL:      fileURL,
		ContentType:  parsedData.ContentType,
		UploadedAt:   time.Now(),
		ElementCount: parsedData.ElementCount,
		DocType:      "knowledge-base.bpt-docs", // Matches Index Type Mapping
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

	// 5. Save to Couchbase
	if err := repositories.SaveDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database save failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ingestion successful",
		"id":      doc.ID,
		"url":     doc.FileURL,
		"count":   doc.ElementCount,
	})
}
