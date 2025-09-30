package middlewares

import (
	"fmt"
	"foglio/v2/src/config"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	tokens    int
	maxTokens int
	refillAt  time.Time
	mutex     sync.RWMutex
}

type RateLimiterConfig struct {
	RequestsPerWindow int
	WindowDuration    time.Duration
	SkipSuccessful    bool
	KeyGenerator      func(*gin.Context) string
	SkipCheck         func(*gin.Context) bool
}

var (
	defaultConfig = RateLimiterConfig{
		RequestsPerWindow: 100,
		WindowDuration:    time.Hour,
		SkipSuccessful:    false,
		KeyGenerator:      defaultKeyGenerator,
		SkipCheck:         defaultSkipCheck,
	}
	limiters        = sync.Map{}
	cleanupInterval = time.Minute * 10
	limiterTTL      = time.Hour * 2
)

func init() {
	go cleanupExpiredLimiters()
}

func defaultKeyGenerator(ctx *gin.Context) string {
	return ctx.ClientIP()
}

func defaultSkipCheck(ctx *gin.Context) bool {
	return false
}

func shouldSkipRateLimit(ctx *gin.Context) bool {
	userInterface, exists := ctx.Get(config.AppConfig.CurrentUser)
	if !exists {
		return false
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		return false
	}

	return user.IsRecruiter || user.IsPremium
}

func RateLimiterMiddleware(configs ...RateLimiterConfig) gin.HandlerFunc {
	config := defaultConfig
	if len(configs) > 0 {
		config = configs[0]
		if config.RequestsPerWindow <= 0 {
			config.RequestsPerWindow = defaultConfig.RequestsPerWindow
		}
		if config.WindowDuration <= 0 {
			config.WindowDuration = defaultConfig.WindowDuration
		}
		if config.KeyGenerator == nil {
			config.KeyGenerator = defaultConfig.KeyGenerator
		}
		if config.SkipCheck == nil {
			config.SkipCheck = defaultConfig.SkipCheck
		}
	}

	return func(ctx *gin.Context) {
		if shouldSkipRateLimit(ctx) || config.SkipCheck(ctx) {
			ctx.Next()
			return
		}

		key := config.KeyGenerator(ctx)
		remaining := getRemainingTokens(key, config)

		if !isAllowed(key, config) {
			ctx.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerWindow))
			ctx.Header("X-RateLimit-Remaining", "0")
			ctx.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(config.WindowDuration).Unix()))
			ctx.Error(lib.NewApiErrror("Rate limit exceeded", http.StatusTooManyRequests))
			ctx.Abort()
			return
		}

		ctx.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerWindow))
		ctx.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		ctx.Next()

		if config.SkipSuccessful && ctx.Writer.Status() >= 200 && ctx.Writer.Status() < 300 {
			refundToken(key, config)
		}
	}
}

func isAllowed(key string, config RateLimiterConfig) bool {
	now := time.Now()
	limiterInterface, exists := limiters.Load(key)

	if !exists {
		limiter := &rateLimiter{
			tokens:    config.RequestsPerWindow - 1,
			maxTokens: config.RequestsPerWindow,
			refillAt:  now.Add(config.WindowDuration),
		}
		limiters.Store(key, limiter)
		return true
	}

	limiter := limiterInterface.(*rateLimiter)
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	if now.After(limiter.refillAt) {
		limiter.tokens = config.RequestsPerWindow - 1
		limiter.refillAt = now.Add(config.WindowDuration)
		limiter.maxTokens = config.RequestsPerWindow
		return true
	}

	if limiter.tokens > 0 {
		limiter.tokens--
		return true
	}

	return false
}

func getRemainingTokens(key string, config RateLimiterConfig) int {
	limiterInterface, exists := limiters.Load(key)
	if !exists {
		return config.RequestsPerWindow
	}

	limiter := limiterInterface.(*rateLimiter)
	limiter.mutex.RLock()
	defer limiter.mutex.RUnlock()

	if time.Now().After(limiter.refillAt) {
		return config.RequestsPerWindow
	}

	return limiter.tokens
}

func refundToken(key string, _ RateLimiterConfig) {
	limiterInterface, exists := limiters.Load(key)
	if !exists {
		return
	}

	limiter := limiterInterface.(*rateLimiter)
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	if limiter.tokens < limiter.maxTokens {
		limiter.tokens++
	}
}

func cleanupExpiredLimiters() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		limiters.Range(func(key, value interface{}) bool {
			limiter := value.(*rateLimiter)
			limiter.mutex.RLock()
			expired := now.Sub(limiter.refillAt) > limiterTTL
			limiter.mutex.RUnlock()

			if expired {
				limiters.Delete(key)
			}
			return true
		})
	}
}

func UserBasedRateLimiter(requestsPerWindow int, windowDuration time.Duration) gin.HandlerFunc {
	config := RateLimiterConfig{
		RequestsPerWindow: requestsPerWindow,
		WindowDuration:    windowDuration,
		KeyGenerator: func(ctx *gin.Context) string {
			if userID, exists := ctx.Get(config.AppConfig.CurrentUserId); exists {
				return "user:" + userID.(string)
			}
			return ctx.ClientIP()
		},
	}
	return RateLimiterMiddleware(config)
}

func RouteBasedRateLimiter(requestsPerWindow int, windowDuration time.Duration) gin.HandlerFunc {
	config := RateLimiterConfig{
		RequestsPerWindow: requestsPerWindow,
		WindowDuration:    windowDuration,
		KeyGenerator: func(ctx *gin.Context) string {
			return ctx.ClientIP() + ":" + ctx.FullPath()
		},
	}
	return RateLimiterMiddleware(config)
}

func RecruiterOrPremiumRateLimiter(standardRequests, premiumRequests int, windowDuration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requests := standardRequests

		userInterface, exists := ctx.Get(config.AppConfig.CurrentUser)
		if exists {
			if user, ok := userInterface.(*models.User); ok && user.IsRecruiter {
				requests = premiumRequests
			}
		}

		config := RateLimiterConfig{
			RequestsPerWindow: requests,
			WindowDuration:    windowDuration,
			KeyGenerator: func(ctx *gin.Context) string {
				if userID, exists := ctx.Get(config.AppConfig.CurrentUserId); exists {
					return "user:" + userID.(string)
				}
				return ctx.ClientIP()
			},
		}

		RateLimiterMiddleware(config)(ctx)
	}
}
