package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRateLimitRouter(maxRequests int, window time.Duration) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimit(maxRequests, window))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func TestRateLimit_RequestsWithinLimit(t *testing.T) {
	maxRequests := 5
	router := setupRateLimitRouter(maxRequests, 1*time.Minute)

	for i := 0; i < maxRequests; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected status %d, got %d", i+1, http.StatusOK, w.Code)
		}
	}
}

func TestRateLimit_RequestsExceedingLimit(t *testing.T) {
	maxRequests := 3
	router := setupRateLimitRouter(maxRequests, 1*time.Minute)

	// Send requests up to the limit -- these should all succeed.
	for i := 0; i < maxRequests; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d (within limit): expected status %d, got %d", i+1, http.StatusOK, w.Code)
		}
	}

	// The next request should be rate limited (429).
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("request exceeding limit: expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}
