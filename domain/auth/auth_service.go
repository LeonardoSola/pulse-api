package auth

import (
	"github.com/LeonardoSola/wagon-api/internal/config"
	"github.com/LeonardoSola/wagon-api/internal/custom_error"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Optionally check signing method
		return []byte(config.Config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return "", custom_error.ErrUnauthorized
	}

	// TODO: Check if token is expired

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", custom_error.ErrUnauthorized
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", custom_error.ErrUnauthorized
	}

	return userID, nil
}
