package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/carantes/go-rate-limiter/lib"
	"github.com/gin-gonic/gin"
)

type server struct {
	e *gin.Engine
}

func rateLimitMiddleware(config map[string]string) gin.HandlerFunc {

	rl, err := lib.NewRateLimiter(config)

	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		// Rate limit check
		userID := c.ClientIP()
		stats, err := rl.Allow(userID)

		if err != nil {
			c.AbortWithStatus(429)
			return
		}

		fmt.Printf("[user]: %s [algorithm]: %s, [capacity] %d, [remaining] %d \n", userID, stats.Algorithm, stats.Capacity, stats.Remaining)

		// Define rate limit headers
		c.Header("X-RateLimit-Algorithm", stats.Algorithm)
		c.Header("X-RateLimit-Limit", strconv.Itoa(stats.Capacity))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(stats.Remaining))

		// TODO: Change this to be the number of seconds to reset, not the time
		c.Header("X-RateLimit-Reset", stats.Reset.Format(time.RFC3339))

		c.Next()
	}
}

func (s *server) Run(addr string) {
	s.e.Run(addr)
}

func NewServer(config map[string]string) *server {
	r := gin.Default()

	// Unlimited requests, have fun
	r.GET("/unlimited", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{"message": "Unlimited, have fun!"})
	})

	r.GET("/limited", rateLimitMiddleware(config), func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{"message": "Limited, dont over use me!"})
	})

	return &server{
		e: r,
	}
}
