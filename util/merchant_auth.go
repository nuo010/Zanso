package util

import (
	"strings"
	"zanso/model"

	"github.com/gin-gonic/gin"
)

const (
	ContextMerchantKey = "currentMerchant"
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

func ExtractMerchantToken(c *gin.Context) string {
	token := strings.TrimSpace(c.GetHeader("Authorization"))
	if token != "" {
		return ParseBearerToken(token)
	}
	return strings.TrimSpace(c.GetHeader("satoken"))
}

func CurrentMerchantID(c *gin.Context) string {
	merchant, ok := GetCurrentMerchant(c)
	if !ok {
		return ""
	}
	return merchant.ID
}

func GetCurrentMerchant(c *gin.Context) (model.Merchant, bool) {
	value, exists := c.Get(ContextMerchantKey)
	if !exists {
		return model.Merchant{}, false
	}
	merchant, ok := value.(model.Merchant)
	return merchant, ok
}
