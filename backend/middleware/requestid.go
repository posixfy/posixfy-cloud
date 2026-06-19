package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// RequestIDHeader is the header used to carry a correlation id across
// the browser -> cloud -> bridge call chain.
const RequestIDHeader = "X-Request-Id"

const requestIDKey = "request_id"

// RequestID echoes an incoming X-Request-Id or generates a new one, stores it
// on the gin context, and sets it on the response so it can be correlated with
// the bridge's logs.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(RequestIDHeader)
		if rid == "" {
			rid = generateID()
		}
		c.Set(requestIDKey, rid)
		c.Header(RequestIDHeader, rid)
		c.Next()
	}
}

// GetRequestID returns the request id stored on the context, or "" if absent.
func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(requestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
