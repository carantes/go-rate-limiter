package lib_test

import (
	"testing"
	"time"

	"github.com/carantes/go-rate-limiter/lib"
	"github.com/carantes/go-rate-limiter/lib/internal/interfaces"
	"github.com/carantes/go-rate-limiter/lib/internal/mocks"
	"github.com/carantes/go-rate-limiter/lib/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testFactorySuite struct {
	suite.Suite
	rlConfig []struct {
		alg    string
		config map[string]string
	}
}

func (s *testFactorySuite) SetupTest() {
	s.rlConfig = []struct {
		alg    string
		config map[string]string
	}{
		{interfaces.TokenBucket.String(), map[string]string{"algorithm": interfaces.TokenBucket.String(), "capacity": "10", "refillRate": "1"}},
		{interfaces.FixedWindow.String(), map[string]string{"algorithm": interfaces.FixedWindow.String(), "capacity": "10", "duration": "5"}},
		{interfaces.SlidingWindowLog.String(), map[string]string{"algorithm": interfaces.SlidingWindowLog.String(), "capacity": "10", "duration": "5"}},
		{interfaces.SlidingWindowCounter.String(), map[string]string{"algorithm": interfaces.SlidingWindowCounter.String(), "capacity": "10", "duration": "5", "weight": "1.0"}},
		// {interfaces.SlidingWindowCounter.String(), map[string]string{"algorithm": interfaces.RedisSlidingWindowCounter.String(), "capacity": "10", "duration": "5", "weight": "1.0"}},
	}
}

func (s *testFactorySuite) TestInvalidAlgorithm() {
	_, err := lib.NewRateLimiter(map[string]string{
		"algorithm": "invalid",
	})

	assert.Error(s.T(), err)
}

func (s *testFactorySuite) TestMissingAlgorithm() {
	_, err := lib.NewRateLimiter(map[string]string{})

	assert.Error(s.T(), err)
}

func (s *testFactorySuite) TestCapacity() {

	// Set current time to future to prevent bucket from refilling
	mocks.Now = func() time.Time {
		return time.Now().Add(time.Second * 5)
	}

	for _, tt := range s.rlConfig {
		s.Run(tt.alg, func() {
			// rate limiter created with config
			rl, err := lib.NewRateLimiter(tt.config)

			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), rl)

			// parse expected capacity from config
			capacity := utils.ParseInt(tt.config["capacity"])

			for i := 0; i < capacity; i++ {
				// allow requests until capacity is reached
				stats, err := rl.Allow("user")

				assert.NoError(s.T(), err)
				assert.Equal(s.T(), capacity, stats.Capacity)
				assert.Equal(s.T(), capacity-i-1, stats.Remaining)
			}

			// no more capacity, throw error
			_, err = rl.Allow("user")
			s.Error(err)
		})
	}
}

func (s *testFactorySuite) TestRefilling() {
	// Set current time to past to allow bucket to refill
	// timer needs to be before the window duration of the algorithm
	mocks.Now = func() time.Time {
		return time.Now().Add(time.Second * 5 * -1)
	}

	for _, tt := range s.rlConfig {
		s.Run(tt.alg, func() {
			// rate limiter created with config
			rl, err := lib.NewRateLimiter(tt.config)

			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), rl)

			// parse expected capacity from config
			capacity := utils.ParseInt(tt.config["capacity"])

			for i := 0; i < capacity; i++ {

				stats, err := rl.Allow("user")

				assert.NoError(s.T(), err)
				assert.Equal(s.T(), capacity, stats.Capacity)

				// remaining tokens should be equal to capacity - 1 (current request)
				assert.Equal(s.T(), capacity-1, stats.Remaining)
			}
		})
	}
}

func TestRateLimiterSuite(t *testing.T) {
	suite.Run(t, new(testFactorySuite))
}
