package handler

import (
	"catalog-gin/config"
	"catalog-gin/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Profile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	var user model.User
	config.DB.Preload("Role.Permissions").First(&user, "id = ?", userID)

	c.JSON(http.StatusOK, user)
}
