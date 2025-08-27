package utils

import (
	"errors"
	"net/http"

	"example.com/m/internal/api/auth"
	"github.com/gin-gonic/gin"
)

func FormResponse(msg string) (int, gin.H) {
	return 200, gin.H{
		"success": true,
		"message": msg,
	}
}

func FormErrResponse(status int, err string) (int, gin.H) {
	return status, gin.H{
		"success": false,
		"error":   err,
	}
}

func FormInvalidRequestResponse() (int, gin.H) {
	return FormErrResponse(http.StatusBadRequest, "invalid request")
}

func FormInternalErrResponse() (int, gin.H) {
	return FormErrResponse(http.StatusInternalServerError, "internal server error")
}

func GetClaims(c *gin.Context) (*auth.Claims, error) {
	claims, ok := c.Get("claims")
	if !ok {
		return nil, errors.New("no claims in context")
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}
	return userClaims, nil
}
