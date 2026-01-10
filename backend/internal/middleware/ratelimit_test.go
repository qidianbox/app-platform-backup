package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter(10, 5) // 10 burst, 5 per second

	// Should allow first 10 requests (burst)
	for i := 0; i < 10; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 11th request should be denied
	if limiter.Allow() {
		t.Error("11th request should be denied")
	}

	// Wait for refill
	time.Sleep(300 * time.Millisecond) // Should refill ~1.5 tokens

	// Should allow 1 more request
	if !limiter.Allow() {
		t.Error("Request after refill should be allowed")
	}
}

func TestIPRateLimiter_GetLimiter(t *testing.T) {
	ipLimiter := NewIPRateLimiter(10, 5)

	// Get limiter for IP1
	limiter1 := ipLimiter.GetLimiter("192.168.1.1")
	if limiter1 == nil {
		t.Error("Limiter should not be nil")
	}

	// Get limiter for same IP should return same instance
	limiter1Again := ipLimiter.GetLimiter("192.168.1.1")
	if limiter1 != limiter1Again {
		t.Error("Should return same limiter for same IP")
	}

	// Get limiter for different IP should return different instance
	limiter2 := ipLimiter.GetLimiter("192.168.1.2")
	if limiter1 == limiter2 {
		t.Error("Should return different limiter for different IP")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router with rate limit middleware
	router := gin.New()
	router.Use(RateLimitMiddleware(5, 2)) // 5 burst, 2 per second
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Make requests
	successCount := 0
	rateLimitedCount := 0

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			successCount++
		} else if w.Code == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	// Should have some successful and some rate limited
	if successCount == 0 {
		t.Error("Should have some successful requests")
	}
	if rateLimitedCount == 0 {
		t.Error("Should have some rate limited requests")
	}
	t.Logf("Success: %d, Rate Limited: %d", successCount, rateLimitedCount)
}

func TestAPIRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router with API rate limit middleware
	router := gin.New()
	router.Use(APIRateLimitMiddleware(3, time.Second)) // 3 requests per second
	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Make 5 requests quickly
	results := make([]int, 5)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		results[i] = w.Code
	}

	// First 3 should succeed, rest should be rate limited
	successCount := 0
	for _, code := range results {
		if code == http.StatusOK {
			successCount++
		}
	}

	if successCount != 3 {
		t.Errorf("Expected 3 successful requests, got %d", successCount)
	}
}

func TestConcurrentRateLimiting(t *testing.T) {
	limiter := NewRateLimiter(100, 50)

	var wg sync.WaitGroup
	var mu sync.Mutex
	allowed := 0

	// Make 200 concurrent requests
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should allow approximately 100 requests (the burst limit)
	if allowed > 110 || allowed < 90 {
		t.Errorf("Expected ~100 allowed requests, got %d", allowed)
	}
	t.Logf("Allowed %d out of 200 concurrent requests", allowed)
}

func TestGlobalRateLimitMiddleware_NotInitialized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset global limiter
	globalIPLimiter = nil

	router := gin.New()
	router.Use(GlobalRateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Should pass through when not initialized
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK when limiter not initialized, got %d", w.Code)
	}
}

func TestGlobalRateLimitMiddleware_Initialized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize global limiter
	InitRateLimiter(5, 2)

	router := gin.New()
	router.Use(GlobalRateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Make requests
	successCount := 0
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			successCount++
		}
	}

	// Should have some rate limiting
	if successCount == 10 {
		t.Error("Expected some requests to be rate limited")
	}
	t.Logf("Success: %d out of 10", successCount)
}
