package util

import (
	"strings"
	"zanso/model"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserKey = "currentUser"
)

func ParseBearerToken(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}
	return authHeader
}

func ExtractUserToken(c *gin.Context) string {
	token := strings.TrimSpace(c.GetHeader("Authorization"))
	if token != "" {
		return ParseBearerToken(token)
	}
	return strings.TrimSpace(c.GetHeader("satoken"))
}

func CurrentUserID(c *gin.Context) string {
	user, ok := GetCurrentUser(c)
	if !ok {
		return ""
	}
	return user.ID
}

func GetCurrentUser(c *gin.Context) (model.User, bool) {
	value, exists := c.Get(ContextUserKey)
	if !exists {
		return model.User{}, false
	}
	user, ok := value.(model.User)
	return user, ok
}
