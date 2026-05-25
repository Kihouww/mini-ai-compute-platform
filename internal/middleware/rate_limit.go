package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type rateLimitResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RateLimitMiddleware(redisClient *redis.Client, limitPerMinute int) gin.HandlerFunc {
	if limitPerMinute <= 0 {
		limitPerMinute = 20
	}

	return func(c *gin.Context) {
		value, exists := c.Get("api_key")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rateLimitResponse{
				Code:    401,
				Message: "missing api key",
				Data:    nil,
			})
			return
		}

		apiKey, ok := value.(string)
		if !ok || apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rateLimitResponse{
				Code:    401,
				Message: "invalid api key",
				Data:    nil,
			})
			return
		}

		minute := time.Now().Format("200601021504")
		key := fmt.Sprintf("rate_limit:%s:%s", apiKey, minute)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, rateLimitResponse{
				Code:    500,
				Message: "rate limit check failed",
				Data:    nil,
			})
			return
		}

		if count == 1 {
			if err := redisClient.Expire(ctx, key, 60*time.Second).Err(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, rateLimitResponse{
					Code:    500,
					Message: "rate limit expire failed",
					Data:    nil,
				})
				return
			}
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limitPerMinute))
		c.Header("X-RateLimit-Used", fmt.Sprintf("%d", count))

		if count > int64(limitPerMinute) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, rateLimitResponse{
				Code:    429,
				Message: "rate limit exceeded",
				Data:    nil,
			})
			return
		}

		c.Next()
	}
}
