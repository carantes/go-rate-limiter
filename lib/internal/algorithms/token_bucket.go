package algorithms

/*
Token Bucket Algorithm
A new bucket is created every time a new user is seen.
The bucket is filled with N tokens and every time a request arrives, a token is removed from the bucket.
If the bucket is empty, the request is declined. The bucket is refilled every second (refill rate) with N tokens.
*/

import (
	"time"

	"github.com/carantes/go-rate-limiter/lib/internal/interfaces"
	"github.com/carantes/go-rate-limiter/lib/internal/mocks"
	"github.com/carantes/go-rate-limiter/lib/internal/utils"
)

// TokenBucketLimiter implements the RateLimiter interface
type tokenBucketLimiter struct {
	usersMap          map[string]*userTokenBucket
	defaultCapacity   int
	defaultRefillRate int
}

// userTokenBucket represents a token bucket for a specific user
type userTokenBucket struct {
	current    int       // current number of available tokens
	capacity   int       // maximum capacity of tokens
	refillRate int       // number of tokens to add per second
	lastRefill time.Time // last time the bucket was refilled
}

type TokenBucketArgs struct {
	Capacity   int
	RefillRate int
}

// Rate Limiter Constructor
func NewTokenBucketLimiter(args TokenBucketArgs) interfaces.RateLimiter {
	return &tokenBucketLimiter{
		usersMap:          make(map[string]*userTokenBucket),
		defaultCapacity:   args.Capacity,
		defaultRefillRate: args.RefillRate,
	}
}

func NewTokenBucketLimiterFromConfig(config map[string]string) (interfaces.RateLimiter, error) {
	capacity, ok := config["capacity"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit capacity"}
	}

	refillRate, ok := config["refillRate"]

	if !ok {
		return nil, &interfaces.RateLimitError{Message: "Missing rate limit refill rate"}
	}

	return NewTokenBucketLimiter(TokenBucketArgs{
		Capacity:   utils.ParseInt(capacity),
		RefillRate: utils.ParseInt(refillRate),
	}), nil
}

func (l *tokenBucketLimiter) Allow(userId string) (interfaces.RateLimiterStats, error) {
	// read user from the map
	bucket := l.usersMap[userId]

	// first request for this user, create a new bucket
	if bucket == nil {
		bucket = &userTokenBucket{
			current:    l.defaultCapacity, // start with full bucket
			capacity:   l.defaultCapacity,
			refillRate: l.defaultRefillRate,
			lastRefill: mocks.Now(),
		}

		l.usersMap[userId] = bucket
	}

	err := bucket.checkTokens()

	if (err) != nil {
		return interfaces.RateLimiterStats{}, err
	}

	return bucket.stats(), nil
}

func (b *userTokenBucket) checkTokens() error {
	// refill the bucket before checking
	b.refill()

	// Not enough tokens to fulfill the request
	if b.current <= 0 {
		return &interfaces.RateLimitError{Message: "Rate limit exceeded"}
	}

	b.current--

	return nil
}

func (b *userTokenBucket) refill() {
	// calculate the number of tokens to add since the last refill
	elapsed := time.Since(b.lastRefill)

	if elapsed.Seconds() <= 0 {
		return
	}

	tokensToAdd := int(elapsed.Seconds()) * b.refillRate

	newCurrent := b.current + tokensToAdd

	// check if the number of tokens exceeds the capacity
	if newCurrent > b.capacity {
		b.current = b.capacity
	} else {
		b.current = newCurrent
	}

	b.lastRefill = mocks.Now()
}

// Return the rate limit stats for the user
func (b *userTokenBucket) stats() interfaces.RateLimiterStats {
	return interfaces.RateLimiterStats{
		Algorithm:   interfaces.TokenBucket.String(),
		Capacity:    b.capacity,
		Remaining:   b.current,
		Reset:       b.lastRefill.Add(time.Duration(b.capacity-b.current) * time.Second),
		CurrentTime: mocks.Now(),
	}
}
