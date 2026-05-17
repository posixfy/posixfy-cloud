package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count    int
	resetAt  time.Time
}

// RateLimit returns a middleware that limits requests per IP.
// maxRequests per window duration.
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	clients := make(map[string]*rateLimitEntry)

	// Background cleanup every window duration
	go func() {
		for {
			time.Sleep(window)
			mu.Lock()
			now := time.Now()
			for ip, entry := range clients {
				if now.After(entry.resetAt) {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		entry, exists := clients[ip]
		now := time.Now()

		if !exists || now.After(entry.resetAt) {
			clients[ip] = &rateLimitEntry{
				count:   1,
				resetAt: now.Add(window),
			}
			mu.Unlock()
			c.Next()
			return
		}

		entry.count++
		if entry.count > maxRequests {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please try again later",
			})
			return
		}
		mu.Unlock()
		c.Next()
	}
}
