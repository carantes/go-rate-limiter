package utils_test

import (
	"testing"
	"time"

	"github.com/carantes/go-rate-limiter/lib/internal/utils"
	"github.com/stretchr/testify/suite"
)

type timeStackSuite struct {
	suite.Suite
	*utils.TimeStack
}

func (s *timeStackSuite) SetupTest() {
	s.TimeStack = utils.NewTimeStack()
}

func (s *timeStackSuite) TestPush() {
	now := time.Now()
	s.Equal(0, s.TimeStack.Size())
	s.TimeStack.Push(now)
	s.Equal(1, s.TimeStack.Size())
	s.Equal(now, s.TimeStack.Peek())

}

func (s *timeStackSuite) TestPop() {
	now := time.Now()

	s.Run("Pop non-empty stack", func() {
		s.Equal(0, s.TimeStack.Size())
		s.TimeStack.Push(now)
		s.Equal(1, s.TimeStack.Size())
		s.Equal(now, s.TimeStack.Pop())
		s.Equal(0, s.TimeStack.Size())
	})

	s.Run("Pop empty stack", func() {
		s.Equal(0, s.TimeStack.Size())
		s.Equal(time.Time{}, s.TimeStack.Pop())
		s.Equal(0, s.TimeStack.Size())
	})
}

func (s *timeStackSuite) TestPeek() {
	now := time.Now()
	s.Equal(0, s.TimeStack.Size())
	s.TimeStack.Push(now)
	s.Equal(1, s.TimeStack.Size())
	s.Equal(now, s.TimeStack.Peek())
	s.Equal(1, s.TimeStack.Size())
}

func (s *timeStackSuite) TestSize() {
	now := time.Now()
	s.Equal(0, s.TimeStack.Size())
	s.TimeStack.Push(now)
	s.Equal(1, s.TimeStack.Size())
	s.TimeStack.Pop()
	s.Equal(0, s.TimeStack.Size())
}

func TestTimeStackSuite(t *testing.T) {
	suite.Run(t, new(timeStackSuite))
}
