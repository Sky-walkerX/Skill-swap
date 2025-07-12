package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Max      int                       // Maximum number of requests
	Duration time.Duration             // Time window
	Message  string                    // Error message when rate limit exceeded
	KeyFunc  func(*gin.Context) string // Function to generate rate limit key
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:      100,         // 100 requests
		Duration: time.Minute, // per minute
		Message:  "Rate limit exceeded. Please try again later.",
		KeyFunc:  func(c *gin.Context) string { return c.ClientIP() },
	}
}

// RateLimit returns a rate limiting middleware
func RateLimit(config ...RateLimitConfig) gin.HandlerFunc {
	cfg := DefaultRateLimitConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Simple in-memory rate limiter
	// For production, consider using Redis
	type client struct {
		count     int
		resetTime time.Time
		mutex     sync.Mutex
	}

	clients := make(map[string]*client)
	mutex := sync.RWMutex{}

	// Cleanup routine
	go func() {
		for {
			time.Sleep(time.Minute)
			mutex.Lock()
			for key, c := range clients {
				c.mutex.Lock()
				if time.Now().After(c.resetTime) {
					delete(clients, key)
				}
				c.mutex.Unlock()
			}
			mutex.Unlock()
		}
	}()

	return func(c *gin.Context) {
		key := cfg.KeyFunc(c)

		mutex.RLock()
		clientData, exists := clients[key]
		mutex.RUnlock()

		if !exists {
			mutex.Lock()
			clientData = &client{
				count:     1,
				resetTime: time.Now().Add(cfg.Duration),
			}
			clients[key] = clientData
			mutex.Unlock()
			c.Next()
			return
		}

		clientData.mutex.Lock()
		defer clientData.mutex.Unlock()

		// Reset if time window has passed
		if time.Now().After(clientData.resetTime) {
			clientData.count = 1
			clientData.resetTime = time.Now().Add(cfg.Duration)
			c.Next()
			return
		}

		// Check if limit exceeded
		if clientData.count >= cfg.Max {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       cfg.Message,
				"retry_after": int(time.Until(clientData.resetTime).Seconds()),
			})
			c.Abort()
			return
		}

		clientData.count++
		c.Next()
	}
}

// AuthRateLimit returns a stricter rate limit for auth endpoints
func AuthRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Max:      5,           // 5 requests
		Duration: time.Minute, // per minute
		Message:  "Too many authentication attempts. Please try again later.",
		KeyFunc:  func(c *gin.Context) string { return c.ClientIP() },
	})
}
