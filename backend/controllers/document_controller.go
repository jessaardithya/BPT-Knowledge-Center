package controllers

import (
	"bpt-knowledge-center/backend/models"
	"bpt-knowledge-center/backend/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDocuments(c *gin.Context) {
	docs, err := repositories.GetAllDocuments()
	if err != nil {
		// Log the error so we can see it in the terminal
		println("Error fetching documents:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	// Return empty list instead of null if no docs found
	if docs == nil {
		docs = []models.Document{}
	}

	c.JSON(http.StatusOK, docs)
}
