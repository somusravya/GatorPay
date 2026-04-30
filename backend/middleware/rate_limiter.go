package middleware

import (
	"net/http"
	"sync"
	"time"

	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// rateLimiter implements a token bucket rate limiter
type rateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int           // tokens per interval
	interval time.Duration // refill interval
	burst    int           // max tokens (bucket capacity)
}

type bucket struct {
	tokens    int
	lastRefill time.Time
}

var limiter *rateLimiter

func init() {
	limiter = &rateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     10,
		interval: time.Second,
		burst:    50,
	}

	// Cleanup goroutine
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			limiter.mu.Lock()
			for key, b := range limiter.buckets {
				if time.Since(b.lastRefill) > 10*time.Minute {
					delete(limiter.buckets, key)
				}
			}
			limiter.mu.Unlock()
		}
	}()
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, exists := rl.buckets[key]
	if !exists {
		rl.buckets[key] = &bucket{
			tokens:    rl.burst - 1,
			lastRefill: time.Now(),
		}
		return true
	}

	// Refill tokens
	elapsed := time.Since(b.lastRefill)
	tokensToAdd := int(elapsed / rl.interval) * rl.rate
	if tokensToAdd > 0 {
		b.tokens += tokensToAdd
		if b.tokens > rl.burst {
			b.tokens = rl.burst
		}
		b.lastRefill = time.Now()
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// RateLimiterMiddleware applies rate limiting per IP address
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !limiter.allow(key) {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}
		c.Next()
	}
}
