package routes

import (
	"bpt-knowledge-center/backend/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the server routes and middleware
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API Group
	api := r.Group("/api")
	{
		api.POST("/documents/upload", controllers.UploadDocument)
		api.GET("/documents", controllers.GetDocuments)
		api.PUT("/documents/:id", controllers.UpdateDocument)
		api.PATCH("/documents/:id/name", controllers.UpdateDocumentName)
		api.DELETE("/documents/:id", controllers.DeleteDocument)

		api.POST("/chat", controllers.HandleChat)
	}

	return r
}
