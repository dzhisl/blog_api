package auth

import (
	"errors"
	"fmt"
	"time"

	"example.com/m/internal/storage/models"
	"example.com/m/internal/types"
	"example.com/m/pkg/config"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"user_role"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User, tokenType types.TokenType) (string, error) {
	var expirationTime time.Duration
	var subject string = string(tokenType)

	switch tokenType {
	case types.TokenAccess:
		expirationTime = time.Minute * 15
	case types.TokenRefresh:
		expirationTime = time.Hour * 24 * 7
	default:
		return "", errors.New("invalid token type")
	}

	claims := &Claims{
		UserID: user.ID,
		Role:   user.UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().SecretKey))
}

func ValidateToken(tokenString string, tokenType types.TokenType, verifyTokenType bool) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(config.GetConfig().SecretKey), nil
	})

	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Subject != string(tokenType) && verifyTokenType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
