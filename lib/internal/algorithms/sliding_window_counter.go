package algorithms

import (
	"time"

	"github.com/carantes/go-rate-limiter/lib/internal/interfaces"
	"github.com/carantes/go-rate-limiter/lib/internal/mocks"
	"github.com/carantes/go-rate-limiter/lib/internal/utils"
)

type slidingWindowCounterLimiter struct {
	userMap               map[string]*userSlidingWindowCounter
	defaultWindowCapacity int
	defaultWindowDuration time.Duration
	currentWindowWeight   float64
}

type userSlidingWindowCounter struct {
	duration                time.Duration
	capacity                int
	currentWindowWeight     float64
	currentWindowStartTime  time.Time
	currentWindowCount      int
	previousWindowStartTime time.Time
	previousWindowCount     int
}

type SlidingWindowCounterArgs struct {
	Capacity int
	Duration time.Duration
	Weight   float64
}

// Rate Limiter Constructor
func NewSlidingWindowCounterLimiter(args SlidingWindowCounterArgs) interfaces.RateLimiter {
	return &slidingWindowCounterLimiter{
		userMap:               make(map[string]*userSlidingWindowCounter),
		defaultWindowCapacity: args.Capacity,
		defaultWindowDuration: args.Duration * time.Second,
		currentWindowWeight:   args.Weight,
	}
}

func NewSlidingWindowCounterLimiterFromConfig(config map[string]string) (interfaces.RateLimiter, error) {
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

	return NewSlidingWindowCounterLimiter(SlidingWindowCounterArgs{
		Capacity: utils.ParseInt(capacity),
		Duration: time.Duration(utils.ParseInt(duration)),
		Weight:   utils.ParseFloat(weight),
	}), nil
}

func (l *slidingWindowCounterLimiter) Allow(user string) (interfaces.RateLimiterStats, error) {
	// read user from the map
	userWindow := l.userMap[user]

	if userWindow == nil {
		userWindow = &userSlidingWindowCounter{
			duration:                l.defaultWindowDuration,
			capacity:                l.defaultWindowCapacity,
			currentWindowWeight:     l.currentWindowWeight,
			currentWindowStartTime:  mocks.Now(),
			currentWindowCount:      0,
			previousWindowStartTime: mocks.Now().Add(-l.defaultWindowDuration),
			previousWindowCount:     0,
		}

		l.userMap[user] = userWindow
	}

	// if user exists, check if there are enough tokens to allow the request
	err := userWindow.checkTokens()

	if (err) != nil {
		return interfaces.RateLimiterStats{}, err
	}

	return userWindow.stats(), nil
}

func (sw *userSlidingWindowCounter) checkTokens() error {
	// check if the current window has expired
	if time.Since(sw.currentWindowStartTime) > sw.duration {
		sw.previousWindowStartTime = sw.currentWindowStartTime
		sw.previousWindowCount = sw.currentWindowCount
		sw.currentWindowStartTime = mocks.Now()
		sw.currentWindowCount = 0
	}

	// check if there are enough tokens to fulfill the request
	if (sw.currentTokens()) >= sw.capacity {
		return &interfaces.RateLimitError{Message: "Rate limit exceeded"}
	}

	// increment the counter
	sw.currentWindowCount++

	return nil
}

func (sw *userSlidingWindowCounter) currentTokens() int {
	previousWindowWeight := 1 - sw.currentWindowWeight

	return int((float64(sw.currentWindowCount)*sw.currentWindowWeight + float64(sw.previousWindowCount)*previousWindowWeight))
}

func (sw *userSlidingWindowCounter) stats() interfaces.RateLimiterStats {
	return interfaces.RateLimiterStats{
		Algorithm: interfaces.SlidingWindowCounter.String(),
		Capacity:  sw.capacity,
		Remaining: sw.capacity - sw.currentTokens(),
		// The reset time is the end of the current window plus the duration of the previous window
		Reset:       sw.currentWindowStartTime.Add(sw.duration * 2),
		CurrentTime: mocks.Now(),
	}
}
