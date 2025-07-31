package router

import (
	"catalog-gin/handler"
	"catalog-gin/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", handler.GetTableCount)
	// Default route returns table count
	r.POST("/login", handler.Login)
	r.POST("/register", handler.Register)
	r.GET("/auth/:provider", handler.BeginOAuth)

	// Route to handle the callback from the OAuth provider after authentication.
	// The provider redirects the user back to this URL.
	// Example: GET /auth/google/callback will process the Google OAuth response.
	r.GET("/auth/:provider/callback", handler.OAuthCallback)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/profile", handler.Profile)
		auth.GET("/admin-only", middleware.RequirePermission("view_admin"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "You have access!"})
		})
	}

	return r
}
