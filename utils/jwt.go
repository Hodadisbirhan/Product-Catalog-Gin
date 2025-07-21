package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateJWT(claims jwt.MapClaims, isAccess bool) (string, error) {

	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	if isAccess {
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		accessTokenString, err := accessToken.SignedString(jwtKey)

		return accessTokenString, err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenString, err := refreshToken.SignedString(jwtKey)

	return refreshTokenString, err
}

func AccessTokenGenerate(userID uuid.UUID) (string, error) {
	return CreateJWT(jwt.MapClaims{
		"exp":     time.Now().Add(time.Minute * time.Duration(15)).Unix(),
		"user_id": userID.String(),
	}, true)
}

func RefreshTokenGenerate(userID uuid.UUID) (string, error) {
	return CreateJWT(jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * time.Duration(48)).Unix(),
		"user_id": userID.String(),
	}, false)
}
