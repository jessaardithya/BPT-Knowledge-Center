package controllers

import (
	"bpt-knowledge-center/backend/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Get Vector from Python Service
	// We ask the Python brain: "What does this question look like as numbers?"
	vector, err := services.GetQueryVector(req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process question (Embedding Service offline?)",
			"details": err.Error(),
		})
		return
	}

	// 2. Search Couchbase (Vector Search)
	// We ask Couchbase: "Find the text chunks most similar to these numbers."
	matches, err := services.SearchSimilarChunks(vector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search knowledge base",
			"details": err.Error(),
		})
		return
	}

	// 3. Formulate Response (Search + Generation)
	// Try to generate a natural answer using Gemini
	answer, err := services.GenerateAnswer(matches, req.Message)

	var responseText string
	if err == nil && answer != "" {
		responseText = answer
	} else {
		// Fallback to raw chunks if LLM fails or no key provided
		if len(matches) == 0 {
			responseText = "I couldn't find any relevant information in your documents."
		} else {
			responseText = "Here is what I found in your documents:\n\n" + strings.Join(matches, "\n\n---\n\n")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"response": responseText,
	})
}
