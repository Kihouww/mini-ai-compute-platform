package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func AuthMiddleware(apiKeys []string) gin.HandlerFunc {
	validKeys := make(map[string]struct{}, len(apiKeys))
	for _, key := range apiKeys {
		validKeys[key] = struct{}{}
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, authResponse{
				Code:    401,
				Message: "missing authorization header",
				Data:    nil,
			})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, authResponse{
				Code:    401,
				Message: "invalid authorization format",
				Data:    nil,
			})
			return
		}

		apiKey := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, authResponse{
				Code:    401,
				Message: "empty api key",
				Data:    nil,
			})
			return
		}

		if _, ok := validKeys[apiKey]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, authResponse{
				Code:    403,
				Message: "invalid api key",
				Data:    nil,
			})
			return
		}

		c.Set("api_key", apiKey)
		c.Next()
	}
}
