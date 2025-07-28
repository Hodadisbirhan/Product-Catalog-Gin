package handler

import (
	"catalog-gin/config"
	"catalog-gin/model"
	"catalog-gin/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
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

type providerCtxKey struct{}

func BeginOAuth(c *gin.Context) {
	// Extract the provider name from the URL parameter (e.g., "google").
	provider := c.Param("provider")

	// Set the provider in the request context. Gothic expects this.
	// We create a new context with the provider value and assign it back to the request.
	ctx := context.WithValue(c.Request.Context(), providerCtxKey{}, provider)
	c.Request = c.Request.WithContext(ctx)

	// Call Gothic's BeginAuthHandler to redirect the user to the OAuth provider's login page.
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func OAuthCallback(c *gin.Context) {
	// Extract the provider name from the URL parameter.
	provider := c.Param("provider")

	// Set the provider in the request context, similar to beginAuth.
	ctx := context.WithValue(c.Request.Context(), providerCtxKey{}, provider)
	c.Request = c.Request.WithContext(ctx)

	log.Printf("Incoming Cookies: %v", c.Request.Cookies())
	if cookie, err := c.Request.Cookie("gothic-session"); err == nil {
		log.Printf("Gothic Session Cookie found: %s", cookie.Value)
	} else {
		log.Printf("Gothic Session Cookie NOT found in request: %v", err)
	}
	// ******************************

	// Complete the user authentication process using Gothic.
	// This exchanges the authorization code for an access token and fetches user data.
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("Error completing user authentication for %s: %v", provider, err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("authentication failed: %w", err))
		return
	}

	// --- Handle User Data ---
	// At this point, the `user` object (type goth.User) contains valuable information:
	// user.Provider: The name of the OAuth provider (e.g., "google")
	// user.Email: User's email address
	// user.Name: User's full name
	// user.FirstName, user.LastName, user.NickName: Additional name components
	// user.AccessToken: OAuth access token (use with caution, typically for API calls to the provider)
	// user.RefreshToken: OAuth refresh token (if "offline_access" scope was requested and granted)
	// user.ExpiresAt: Expiry time of the access token

	// In a real application, you would:
	// 1. Check if this user (identified by user.Provider and user.UserID or user.Email) exists in your database.
	// 2. If not, create a new user entry.
	// 3. Update existing user information if necessary.
	// 4. Securely store sensitive tokens (like RefreshToken) in your database, encrypted.
	// 5. Generate your own session token (e.g., a JWT or a server-side session ID) for the client.
	//    This token would be used for subsequent authenticated requests to your API.

	// For demonstration purposes, we'll simply return the user information as JSON.
	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Printf("Error marshalling user data: %v", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to process user data: %w", err))
		return
	}

	// Respond with the user's data. In a production app, you'd likely redirect
	// to a client-side route with a generated session token.
	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("Successfully logged in with %s!", provider),
		"user_info": string(userJSON),
		// Example: "your_app_session_token": "some_jwt_or_session_id",
	})
	// Alternatively, redirect to a success page or dashboard:
	// c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}
