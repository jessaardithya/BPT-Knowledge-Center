package controllers

import (
	"bpt-knowledge-center/backend/models"
	"bpt-knowledge-center/backend/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateDocRequest struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}

type UpdateNameRequest struct {
	DisplayName string `json:"display_name"`
}

func UpdateDocument(c *gin.Context) {
	id := c.Param("id")

	var req UpdateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	err := repositories.UpdateDocumentMetadata(id, req.Category, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// UpdateDocumentName updates only the display name of a document
func UpdateDocumentName(c *gin.Context) {
	id := c.Param("id")

	var req UpdateNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	if req.DisplayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Display name cannot be empty"})
		return
	}

	err := repositories.UpdateDocumentName(id, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document name"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document name updated successfully"})
}

func DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	err := repositories.DeleteDocument(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
}

func GetDocuments(c *gin.Context) {
	docs, err := repositories.GetAllDocuments()
	if err != nil {
		println("Error fetching documents:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	if docs == nil {
		docs = []models.Document{}
	}

	c.JSON(http.StatusOK, docs)
}
