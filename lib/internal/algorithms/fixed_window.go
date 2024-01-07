package algorithms

import (
	"time"

	"github.com/carantes/go-rate-limiter/lib/internal/interfaces"
	"github.com/carantes/go-rate-limiter/lib/internal/mocks"
	"github.com/carantes/go-rate-limiter/lib/internal/utils"
)

/*
Fixed Window Algorithm
Check the elapsed time between two requests. If the elapsed time is bigger than the window size, then reset the window counter.
Otherwise, increment the counter and check if it is bigger than the maximum number of requests allowed.
The downside of this algorithm is that it allows bursts of requests at the end and beginning of each window.
*/

// FixedWindowLimiter implements the RateLimiter interface
type fixedWindowLimiter struct {
	usersMap              map[string]*userFixedWindow
	defaultWindowCapacity int
	defaultWindowDuration time.Duration
}

// userFixedWindow represents a fixed window for a specific user
type userFixedWindow struct {
	startTime time.Time     // start time of the window
	duration  time.Duration // window size
	capacity  int           // maximum number of requests allowed
	current   int           // current number of requests
}

type FixedWindowArgs struct {
	Capacity int
	Duration time.Duration
}

// Rate Limiter Constructor
func NewFixedWindowLimiter(args FixedWindowArgs) interfaces.RateLimiter {
	return &fixedWindowLimiter{
		usersMap:              make(map[string]*userFixedWindow),
		defaultWindowCapacity: args.Capacity,
		defaultWindowDuration: args.Duration * time.Second,
	}
}

func NewFixedWindowLimiterFromConfig(config map[string]string) (interfaces.RateLimiter, error) {
	capacity, ok := config["capacity"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit capacity"}
	}

	duration, ok := config["duration"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit duration"}
	}

	return NewFixedWindowLimiter(FixedWindowArgs{
		Capacity: utils.ParseInt(capacity),
		Duration: time.Duration(utils.ParseInt(duration)),
	}), nil
}

func (l *fixedWindowLimiter) Allow(user string) (interfaces.RateLimiterStats, error) {
	// read user from the map
	userWindow := l.usersMap[user]

	if userWindow == nil {
		userWindow = &userFixedWindow{
			startTime: mocks.Now(),
			duration:  l.defaultWindowDuration,
			capacity:  l.defaultWindowCapacity,
			current:   0,
		}

		l.usersMap[user] = userWindow
	}

	// if user exists, check if there are enough tokens to allow the request
	err := userWindow.checkTokens()

	if (err) != nil {
		return interfaces.RateLimiterStats{}, err
	}

	return userWindow.stats(), nil
}

func (fw *userFixedWindow) checkTokens() error {
	// check if the window has expired
	if time.Since(fw.startTime) > fw.duration {
		fw.startTime = mocks.Now()
		fw.current = 0
	}

	// check if there are enough tokens to fulfill the request
	if fw.current >= fw.capacity {
		return &interfaces.RateLimitError{Message: "Rate limit exceeded"}
	}

	// increment the counter
	fw.current++

	return nil
}

// Return the rate limit stats for the user
func (fw *userFixedWindow) stats() interfaces.RateLimiterStats {
	return interfaces.RateLimiterStats{
		Algorithm:   interfaces.FixedWindow.String(),
		Capacity:    fw.capacity,
		Remaining:   fw.capacity - fw.current,
		Reset:       fw.startTime.Add(fw.duration),
		CurrentTime: mocks.Now(),
	}
}
