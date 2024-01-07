package lib

import (
	"github.com/carantes/go-rate-limiter/lib/internal/algorithms"
	"github.com/carantes/go-rate-limiter/lib/internal/interfaces"
)

// Rate limiter factory
func NewRateLimiter(config map[string]string) (interfaces.RateLimiter, error) {

	var alg, ok = interfaces.ParseAlgorithm(config["algorithm"])

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit algorithm"}
	}

	switch alg {
	case interfaces.TokenBucket:
		return algorithms.NewTokenBucketLimiterFromConfig(config)
	case interfaces.FixedWindow:
		return algorithms.NewFixedWindowLimiterFromConfig(config)
	case interfaces.SlidingWindowLog:
		return algorithms.NewSlidingWindowLogLimiterFromConfig(config)
	case interfaces.SlidingWindowCounter:
		return algorithms.NewSlidingWindowCounterLimiterFromConfig(config)
	case interfaces.RedisSlidingWindowCounter:
		return algorithms.NewRedisSlidingWindowCounterLimiterFromConfig(config)
	default:
		return nil, &interfaces.RateLimitError{Message: "Invalid rate limit algorithm"}
	}
}
