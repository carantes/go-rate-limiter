package algorithms

import (
	"time"

	"github.com/carantes/go-rate-limiter/lib/internal/interfaces"
	"github.com/carantes/go-rate-limiter/lib/internal/mocks"
	"github.com/carantes/go-rate-limiter/lib/internal/utils"
)

/*
Sliding Window Algorithm
Log the timestamp of each request. When a new request arrives, remove all the timestamps that are older than the window size.
If the number of remaining timestamps is bigger than the maximum number of requests allowed, then decline the request.
The disadvantage of this algorithm is that it requires a lot of memory to store all the timestamps.
*/

// SlidingWindowLimiter implements the RateLimiter interface
type slidingWindowLogLimiter struct {
	usersMap              map[string]*userSlidingWindow
	defaultWindowCapacity int
	defaultWindowDuration time.Duration
}

// userSlidingWindow represents a sliding window for a specific user
type userSlidingWindow struct {
	duration     time.Duration    // window size
	capacity     int              // maximum number of requests allowed
	requestStack *utils.TimeStack // timestamps of the requests
}

type SlidingWindowLogArgs struct {
	Capacity int
	Duration time.Duration
}

// Rate Limiter Constructor
func NewSlidingWindowLogLimiter(args SlidingWindowLogArgs) interfaces.RateLimiter {
	return &slidingWindowLogLimiter{
		usersMap:              make(map[string]*userSlidingWindow),
		defaultWindowCapacity: args.Capacity,
		defaultWindowDuration: args.Duration * time.Second,
	}
}

func NewSlidingWindowLogLimiterFromConfig(config map[string]string) (interfaces.RateLimiter, error) {
	capacity, ok := config["capacity"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit capacity"}
	}

	duration, ok := config["duration"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit duration"}
	}

	return NewSlidingWindowLogLimiter(SlidingWindowLogArgs{
		Capacity: utils.ParseInt(capacity),
		Duration: time.Duration(utils.ParseInt(duration)),
	}), nil
}

func (l *slidingWindowLogLimiter) Allow(user string) (interfaces.RateLimiterStats, error) {
	// read user from the map
	userWindow := l.usersMap[user]

	if userWindow == nil {
		userWindow = &userSlidingWindow{
			duration:     l.defaultWindowDuration,
			capacity:     l.defaultWindowCapacity,
			requestStack: utils.NewTimeStack(),
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

func (sw *userSlidingWindow) checkTokens() error {
	// inline remove requests that are older than the window size
	for sw.requestStack.Size() > 0 {
		if time.Since(sw.requestStack.Peek()) > sw.duration {
			sw.requestStack.Pop()
		} else {
			break
		}
	}

	// add current request timestamp
	sw.requestStack.Push(mocks.Now())

	// check if there are enough tokens to fulfill the request
	if sw.requestStack.Size() > sw.capacity {
		return &interfaces.RateLimitError{Message: "Rate limit exceeded"}
	}

	return nil
}

func (sw *userSlidingWindow) stats() interfaces.RateLimiterStats {
	return interfaces.RateLimiterStats{
		Algorithm:   interfaces.SlidingWindowLog.String(),
		Capacity:    sw.capacity,
		Remaining:   sw.capacity - sw.requestStack.Size(),
		Reset:       mocks.Now().Add(sw.duration),
		CurrentTime: mocks.Now(),
	}
}
