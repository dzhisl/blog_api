package utils

import (
	"errors"
	"net/http"

	"example.com/m/internal/api/auth"
	"example.com/m/internal/types"
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

// CompareRoles compares two roles based on their hierarchy defined in rolesMap.
// It returns:
//
//	-1 if role1 is lower than role2,
//	 0 if both roles are equal,
//	 1 if role1 is higher than role2.
//
// If either role is not found in rolesMap, it returns 0.
func CompareRoles(role1, role2 types.Role) int {
	var idx1, idx2 int = -1, -1

	for k, v := range types.RolesMap {
		if v == role1 {
			idx1 = k
		}
		if v == role2 {
			idx2 = k
		}
	}

	if idx1 == -1 || idx2 == -1 {
		// unknown role → consider them equal
		return 0
	}

	switch {
	case idx1 < idx2:
		return -1
	case idx1 > idx2:
		return 1
	default:
		return 0
	}
}
