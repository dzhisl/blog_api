package utils

import (
	"net/http"

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
