package middleware

import (
	"NoRethink/internal/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RateLimitConfig struct {
	RedisClient *redis.Client
	KeyPrefix   string
	Limit       redis_rate.Limit
}

func newLimiter(cfg RateLimitConfig) *redis_rate.Limiter {
	return redis_rate.NewLimiter(cfg.RedisClient)
}

func setHeaders(c *gin.Context, cfg RateLimitConfig, res *redis_rate.Result) {
	c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.Limit.Rate))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
}

func abort429(c *gin.Context, cfg RateLimitConfig, res *redis_rate.Result) {
	setHeaders(c, cfg, res)
	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"message": "rate limit exceeded.",
	})
}

func handleRedisError(c *gin.Context, key string, err error) {
	log.Error().Err(err).Str("key", key).Msg("rate limiter: redis error")
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error":   "rate_limiter_unavailable",
		"message": "Rate limiter is temporarily unavailable.",
	})
}

func IPRateLimit(cfg RateLimitConfig) gin.HandlerFunc {
	limiter := newLimiter(cfg)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "missing_client_ip",
				"message": "Could not determine client IP address.",
			})
			return
		}

		key := fmt.Sprintf("ip:%s:%s", cfg.KeyPrefix, ip)

		res, err := limiter.Allow(c.Request.Context(), key, cfg.Limit)
		if err != nil {
			handleRedisError(c, key, err)
			return
		}

		setHeaders(c, cfg, res)

		// redis_rate sets Allowed = 0 when the request is denied.
		if res.Allowed == 0 {
			abort429(c, cfg, res)
			return
		}

		c.Next()
	}
}

func UserRateLimit(cfg RateLimitConfig) gin.HandlerFunc {
	limiter := newLimiter(cfg)

	return func(c *gin.Context) {
		raw, exists := c.Get(util.UserIDKey)
		if !exists || raw == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "missing_user_id",
				"message": "Authenticated user ID not found in request context.",
			})
			return
		}

		userID, ok := raw.(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_user_id",
				"message": "User ID in context is not a valid non-empty string.",
			})
			return
		}

		key := fmt.Sprintf("user:%s:%s", cfg.KeyPrefix, userID)

		res, err := limiter.Allow(c.Request.Context(), key, cfg.Limit)
		if err != nil {
			handleRedisError(c, key, err)
			return
		}

		setHeaders(c, cfg, res)

		if res.Allowed == 0 {
			abort429(c, cfg, res)
			return
		}

		c.Next()
	}
}
