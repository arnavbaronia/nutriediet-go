package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitAuth creates rate limiter for authentication endpoints
// Limits to 5 requests per minute to prevent brute force attacks
func RateLimitAuth() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	
	middleware := mgin.NewMiddleware(instance)
	return middleware
}

// RateLimitAPI creates rate limiter for general API endpoints
// Limits to 100 requests per minute for normal API usage
func RateLimitAPI() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	
	middleware := mgin.NewMiddleware(instance)
	return middleware
}

// RateLimitStrict creates stricter rate limiter for sensitive endpoints
// Limits to 3 requests per minute for extra sensitive operations
func RateLimitStrict() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  3,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	
	middleware := mgin.NewMiddleware(instance)
	return middleware
}

