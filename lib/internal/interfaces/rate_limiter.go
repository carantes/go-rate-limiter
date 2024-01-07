package interfaces

import "time"

// RateLimiter is an interface that defines the methods that a rate limiter should implement
type RateLimiter interface {
	//check if a request is allowed, return user stats or error
	Allow(user string) (RateLimiterStats, error)
}

// RateLimiterStats represents the stats of a rate limiter for a specific user
type RateLimiterStats struct {
	Algorithm   string
	Capacity    int
	Remaining   int
	Reset       time.Time
	CurrentTime time.Time
}
