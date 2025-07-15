package middleware

import (
	"catalog-gin/config"
	"catalog-gin/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequirePermission(permissionName string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("user_id").(uuid.UUID)
        var user model.User
        config.DB.Preload("Role.Permissions").First(&user, "id = ?", userID)

        for _, perm := range user.Role.Permissions {
            if perm.Name == permissionName {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
        c.Abort()
    }
}
