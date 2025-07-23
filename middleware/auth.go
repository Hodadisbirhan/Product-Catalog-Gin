package middleware

import (
	"catalog-gin/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		secretKey := []byte(os.Getenv("JWT_SECRET"))

		accessTokenString, err := c.Cookie("access")
		if err != nil {
			if err == http.ErrNoCookie {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token provided"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving access token"})
			}
			c.Abort()
			return
		}

		accessToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err == nil && accessToken.Valid {

			if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && claims["user_id"] != nil {
				if userIDStr, ok := claims["user_id"].(string); ok {
					parsedUserID, uuidParseErr := uuid.Parse(userIDStr)
					if uuidParseErr == nil {
						c.Set("user_id", parsedUserID)
						c.Next()
						return
					} else {

						c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed user_id in access token"})
						c.Abort()
						return
					}
				} else {

					c.JSON(http.StatusBadRequest, gin.H{"error": "User ID claim missing or invalid type in access token"})
					c.Abort()
					return
				}
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token claims or missing user_id"})
			c.Abort()
			return
		}

		fmt.Println("Access token failed validation, attempting refresh...")
		refreshTokenString, err := c.Cookie("refresh")
		if err != nil {
			if err == http.ErrNoCookie {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token invalid/expired and no refresh token provided"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving refresh token"})
			}
			c.Abort()
			return
		}

		refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !refreshToken.Valid {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token invalid or expired"})
			c.Abort()
			return
		}

		if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && claims["user_id"] != nil {
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User ID claim missing or not a string in refresh token"})
				c.Abort()
				return
			}

			parsedUserID, uuidParseErr := uuid.Parse(userIDStr)
			if uuidParseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed UUID in refresh token"})
				c.Abort()
				return
			}

			newAccessTokenString, err := utils.AccessTokenGenerate(parsedUserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
				c.Abort()
				return
			}

			c.SetCookie(
				"access",
				newAccessTokenString,
				900,
				"/",
				"localhost",
				false,
				true,
			)

			c.SetCookie(
				"refresh",
				refreshTokenString,
				172800,
				"/",
				"localhost",
				false,
				true,
			)

			c.Set("user_id", parsedUserID)
			c.Next()
			return
		} else {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token claims or missing user_id"})
			c.Abort()
			return
		}
	}
}
