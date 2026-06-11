package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*ipLimiter
	rate     rate.Limit
	burst    int
	cleanup  time.Duration
}

func NewRateLimiter(r rate.Limit, burst int, cleanup time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ipLimiter),
		rate:    r,
		burst:   burst,
		cleanup: cleanup,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	for {
		time.Sleep(rl.cleanup)
		rl.mu.Lock()
		for ip, l := range rl.clients {
			if time.Since(l.lastSeen) > rl.cleanup {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	l, ok := rl.clients[ip]
	if !ok {
		l = &ipLimiter{limiter: rate.NewLimiter(rl.rate, rl.burst)}
		rl.clients[ip] = l
	}
	l.lastSeen = time.Now()
	return l.limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}
