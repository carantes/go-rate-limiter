package algorithms

import (
	"time"

	"github.com/carantes/go-rate-limiter/lib/internal/interfaces"
	"github.com/carantes/go-rate-limiter/lib/internal/mocks"
	"github.com/carantes/go-rate-limiter/lib/internal/utils"
)

type redisSlidingWindowCounterLimiter struct {
	redisClient           *utils.RedisClient
	defaultWindowCapacity int
	defaultWindowDuration time.Duration
	currentWindowWeight   float64
}

type redisUserSlidingWindowCounter struct {
	Duration                time.Duration
	Capacity                int
	CurrentWindowWeight     float64
	CurrentWindowStartTime  time.Time
	CurrentWindowCount      int
	PreviousWindowStartTime time.Time
	PreviousWindowCount     int
}

type RedisSlidingWindowCounterArgs struct {
	RedisURL string
	Capacity int
	Duration time.Duration
	Weight   float64
}

// Rate Limiter Constructor
func NewRedisSlidingWindowCounterLimiter(args RedisSlidingWindowCounterArgs) interfaces.RateLimiter {
	client := utils.NewRedisClient(args.RedisURL)

	return &redisSlidingWindowCounterLimiter{
		redisClient:           client,
		defaultWindowCapacity: args.Capacity,
		defaultWindowDuration: args.Duration * time.Second,
		currentWindowWeight:   args.Weight,
	}
}

func NewRedisSlidingWindowCounterLimiterFromConfig(config map[string]string) (interfaces.RateLimiter, error) {
	capacity, ok := config["capacity"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit capacity"}
	}

	duration, ok := config["duration"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit duration"}
	}

	weight, ok := config["weight"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit weight"}
	}

	redisURL, ok := config["redisURL"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing redis URL"}
	}

	return NewRedisSlidingWindowCounterLimiter(RedisSlidingWindowCounterArgs{
		RedisURL: redisURL,
		Capacity: utils.ParseInt(capacity),
		Duration: time.Duration(utils.ParseInt(duration)),
		Weight:   utils.ParseFloat(weight),
	}), nil
}

func (l *redisSlidingWindowCounterLimiter) Allow(user string) (interfaces.RateLimiterStats, error) {
	// read user from the map
	var userWindow *redisUserSlidingWindowCounter

	redisTTL := l.defaultWindowDuration * 2

	l.redisClient.Get(user, &userWindow)

	if userWindow == nil {
		userWindow = &redisUserSlidingWindowCounter{
			Duration:                l.defaultWindowDuration,
			Capacity:                l.defaultWindowCapacity,
			CurrentWindowWeight:     l.currentWindowWeight,
			CurrentWindowStartTime:  mocks.Now(),
			CurrentWindowCount:      0,
			PreviousWindowStartTime: mocks.Now().Add(-l.defaultWindowDuration),
			PreviousWindowCount:     0,
		}

		// save user window
		l.redisClient.Set(user, userWindow, redisTTL)
	}

	// if user exists, check if there are enough tokens to allow the request
	err := userWindow.checkTokens()

	// update user window
	l.redisClient.Set(user, userWindow, redisTTL)

	if (err) != nil {
		return interfaces.RateLimiterStats{}, err
	}

	return userWindow.stats(), nil
}

func (sw *redisUserSlidingWindowCounter) checkTokens() error {
	// check if the current window has expired
	if time.Since(sw.CurrentWindowStartTime) > sw.Duration {
		sw.PreviousWindowStartTime = sw.CurrentWindowStartTime
		sw.PreviousWindowCount = sw.CurrentWindowCount
		sw.CurrentWindowStartTime = mocks.Now()
		sw.CurrentWindowCount = 0
	}

	// check if there are enough tokens to fulfill the request
	if (sw.currentTokens()) >= sw.Capacity {
		return &interfaces.RateLimitError{Message: "Rate limit exceeded"}
	}

	// increment the counter
	sw.CurrentWindowCount++

	return nil
}

func (sw *redisUserSlidingWindowCounter) currentTokens() int {
	previousWindowWeight := 1 - sw.CurrentWindowWeight

	return int((float64(sw.CurrentWindowCount)*sw.CurrentWindowWeight + float64(sw.PreviousWindowCount)*previousWindowWeight))
}

func (sw *redisUserSlidingWindowCounter) stats() interfaces.RateLimiterStats {
	return interfaces.RateLimiterStats{
		Algorithm: interfaces.RedisSlidingWindowCounter.String(),
		Capacity:  sw.Capacity,
		Remaining: sw.Capacity - sw.currentTokens(),
		// The reset time is the end of the current window plus the duration of the previous window
		Reset:       sw.CurrentWindowStartTime.Add(sw.Duration * 2),
		CurrentTime: mocks.Now(),
	}
}
