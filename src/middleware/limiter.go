package middleware

import (
	"RestGoTest/src/config"
	"RestGoTest/src/helper"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func LimitByRequest() gin.HandlerFunc {
	config := config.GetConfig()
	lmt := tollbooth.NewLimiter(config.Server.RateLimitCount, nil)
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				helper.GenerateBaseResponseWithError(nil, false, helper.LimiterError, err))
			return
		} else {
			c.Next()
		}
	}
}

// visitor bundles a rate.Limiter with its last‑seen timestamp.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter holds the per‑IP limiters and cleans them periodically.
type IPRateLimiter struct {
	mu        sync.RWMutex
	visitors  map[string]*visitor
	r         rate.Limit
	b         int
	ttl       time.Duration // how long to keep an idle IP before eviction
	cleanupTk *time.Ticker
	quit      chan struct{}
}

// NewIPRateLimiter creates a limiter that allows 'b' events every 'r' seconds per IP
// and drops idle IPs after ttl.
func NewIPRateLimiter(r rate.Limit, b int, ttl time.Duration) *IPRateLimiter {
	// اگر ttl مقدار نداشت یا صفر بود، مقدار پیش‌فرض بده (مثلاً 1 دقیقه)
	if ttl <= 0 {
		ttl = time.Minute
	}

	l := &IPRateLimiter{
		visitors:  make(map[string]*visitor),
		r:         r,
		b:         b,
		ttl:       ttl,
		cleanupTk: time.NewTicker(ttl),
		quit:      make(chan struct{}),
	}
	go l.cleanupLoop()
	return l
}

// GetLimiter returns the limiter for the given IP, creating one if necessary.
func (l *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, exists := l.visitors[ip]
	if !exists {
		v = &visitor{
			limiter:  rate.NewLimiter(l.r, l.b),
			lastSeen: time.Now(),
		}
		l.visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

// Stop terminates the background cleanup goroutine.
func (l *IPRateLimiter) Stop() {
	close(l.quit)
}

// cleanupLoop evicts idle IPs periodically.
func (l *IPRateLimiter) cleanupLoop() {
	for {
		select {
		case <-l.cleanupTk.C:
			l.evict()
		case <-l.quit:
			l.cleanupTk.Stop()
			return
		}
	}
}

// evict removes entries not seen for > ttl.
func (l *IPRateLimiter) evict() {
	l.mu.Lock()
	defer l.mu.Unlock()
	cutoff := time.Now().Add(-l.ttl)
	for ip, v := range l.visitors {
		if v.lastSeen.Before(cutoff) {
			delete(l.visitors, ip)
		}
	}
}

// ------------------- Gin middleware ------------------- //

// OtpLimiter middleware enforces per‑IP rate limits for OTP endpoints.
// It uses a shared IPRateLimiter instance with automatic cleanup.
func OtpLimiter(cfg *config.Config) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Every(cfg.Otp.Limiter*time.Second), cfg.Otp.Burst, cfg.Otp.Ttl*time.Minute) // keep idle IPs 15m

	return func(c *gin.Context) {
		ip := clientIP(c.Request.RemoteAddr)
		if !limiter.GetLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				helper.GenerateBaseResponseWithError(nil, false, helper.OtpLimiterError, errors.New("not allowed")))
			return
		}
		c.Next()
	}
}

// clientIP extracts the IP from RemoteAddr safely.
func clientIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr // already just IP
	}
	return ip
}
