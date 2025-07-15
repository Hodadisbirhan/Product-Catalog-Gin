package router

import (
	"catalog-gin/handler"
	"catalog-gin/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
	r.GET("/", handler.GetTableCount) 
    r.POST("/login", handler.Login)
 // Default route returns table count

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
