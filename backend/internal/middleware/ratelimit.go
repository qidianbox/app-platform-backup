package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 令牌桶限流器
type RateLimiter struct {
	mu           sync.Mutex
	tokens       float64
	maxTokens    float64
	refillRate   float64 // 每秒补充的令牌数
	lastRefill   time.Time
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(maxTokens float64, refillRate float64) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许请求
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.tokens += elapsed * r.refillRate
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
	r.lastRefill = now

	// 检查是否有足够的令牌
	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}

// IPRateLimiter 基于IP的限流器
type IPRateLimiter struct {
	mu        sync.RWMutex
	limiters  map[string]*RateLimiter
	maxTokens float64
	refillRate float64
	cleanupInterval time.Duration
}

// NewIPRateLimiter 创建基于IP的限流器
func NewIPRateLimiter(maxTokens, refillRate float64) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters:   make(map[string]*RateLimiter),
		maxTokens:  maxTokens,
		refillRate: refillRate,
		cleanupInterval: 10 * time.Minute,
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter
}

// GetLimiter 获取指定IP的限流器
func (l *IPRateLimiter) GetLimiter(ip string) *RateLimiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.limiters[ip]
	if !exists {
		limiter = NewRateLimiter(l.maxTokens, l.refillRate)
		l.limiters[ip] = limiter
	}

	return limiter
}

// cleanup 定期清理过期的限流器
func (l *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(l.cleanupInterval)
	for range ticker.C {
		l.mu.Lock()
		// 清理超过1小时未使用的限流器
		threshold := time.Now().Add(-1 * time.Hour)
		for ip, limiter := range l.limiters {
			if limiter.lastRefill.Before(threshold) {
				delete(l.limiters, ip)
			}
		}
		l.mu.Unlock()
	}
}

// 全局IP限流器
var globalIPLimiter *IPRateLimiter

// InitRateLimiter 初始化全局限流器
func InitRateLimiter(maxTokens, refillRate float64) {
	globalIPLimiter = NewIPRateLimiter(maxTokens, refillRate)
}

// RateLimitMiddleware 限流中间件
// maxTokens: 最大令牌数（突发请求上限）
// refillRate: 每秒补充的令牌数（持续请求速率）
func RateLimitMiddleware(maxTokens, refillRate float64) gin.HandlerFunc {
	limiter := NewIPRateLimiter(maxTokens, refillRate)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GlobalRateLimitMiddleware 使用全局限流器的中间件
func GlobalRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if globalIPLimiter == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		if !globalIPLimiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// APIRateLimitMiddleware 针对特定API的限流中间件
func APIRateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	type requestInfo struct {
		count     int
		resetTime time.Time
	}

	var mu sync.Mutex
	requests := make(map[string]*requestInfo)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path
		key := ip + ":" + path

		mu.Lock()
		info, exists := requests[key]
		now := time.Now()

		if !exists || now.After(info.resetTime) {
			requests[key] = &requestInfo{
				count:     1,
				resetTime: now.Add(window),
			}
			info = requests[key]
			mu.Unlock()
			// 添加限流响应头
			c.Header("X-RateLimit-Limit", intToStr(maxRequests))
			c.Header("X-RateLimit-Remaining", intToStr(maxRequests-1))
			c.Header("X-RateLimit-Reset", intToStr(int(info.resetTime.Unix())))
			c.Next()
			return
		}

		remaining := maxRequests - info.count - 1
		if remaining < 0 {
			remaining = 0
		}

		if info.count >= maxRequests {
			retryAfter := int(info.resetTime.Sub(now).Seconds())
			mu.Unlock()
			// 添加限流响应头
			c.Header("X-RateLimit-Limit", intToStr(maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", intToStr(int(info.resetTime.Unix())))
			c.Header("Retry-After", intToStr(retryAfter))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":        429,
				"message":     "该接口请求过于频繁，请稍后再试",
				"retry_after": retryAfter,
			})
			c.Abort()
			return
		}

		info.count++
		mu.Unlock()
		// 添加限流响应头
		c.Header("X-RateLimit-Limit", intToStr(maxRequests))
		c.Header("X-RateLimit-Remaining", intToStr(remaining))
		c.Header("X-RateLimit-Reset", intToStr(int(info.resetTime.Unix())))
		c.Next()
	}
}

// intToStr 将int转换为字符串
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	negative := n < 0
	if negative {
		n = -n
	}
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	if negative {
		result = append([]byte{'-'}, result...)
	}
	return string(result)
}
