package handler

import (
	"catalog-gin/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTableCount(c *gin.Context) {
	var count int64
	// Query the number of tables in the public schema
	err := config.DB.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`).Scan(&count).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count tables"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"table_count": count,
	})
}
