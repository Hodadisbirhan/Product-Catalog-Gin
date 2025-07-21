package handler

import (
	"catalog-gin/config"
	"catalog-gin/model"
	"catalog-gin/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	LoginInput
	Username string `json:"username" binding:"required,min=3,max=100"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	result := config.DB.Preload("Role.Permissions").Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := utils.AccessTokenGenerate(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Access token generation failed"})
		return
	}

	refreshToken, err := utils.RefreshTokenGenerate(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh token generation failed"})
		return
	}

	c.SetCookie(
		"access",
		accessToken,
		900, // 15 minutes
		"/",
		"localhost",
		false,
		true,
	)

	c.SetCookie(
		"refresh",
		refreshToken,
		172800, // 48 hours
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user model.User
	resultUser := config.DB.Where("email = ?", input.Email).First(&user)
	if resultUser.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	input.Password = hashedPassword

	var defaultRole model.Role
	role := config.DB.Where("Name = ?", "user").First(&defaultRole)
	if role.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Role not found"})
		return
	}

	userRecord := model.User{
		Email:    input.Email,
		Username: input.Username,
		Password: input.Password,
		RoleID:   defaultRole.ID,
	}

	if err := config.DB.Create(&userRecord).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
		} else {
			log.Printf("Error creating user in database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	accessToken, err := utils.AccessTokenGenerate(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Access token generation failed"})
		return
	}

	refreshToken, err := utils.RefreshTokenGenerate(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh token generation failed"})
		return
	}

	c.SetCookie(
		"access",
		accessToken,
		900, // 15 minutes
		"/",
		"localhost",
		false,
		true,
	)

	c.SetCookie(
		"refresh",
		refreshToken,
		172800, // 48 hours
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
