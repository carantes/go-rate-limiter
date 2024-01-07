package utils_test

import (
	"testing"

	"github.com/carantes/go-rate-limiter/lib/internal/utils"
	"github.com/stretchr/testify/suite"
)

type parserSuite struct {
	suite.Suite
}

func (s *parserSuite) TestParseInt() {
	s.Run("Parse valid int", func() {
		s.Equal(1234, utils.ParseInt("1234"))
	})

	s.Run("Parse invalid int", func() {
		s.Equal(0, utils.ParseInt("abc"))
	})
}

func (s *parserSuite) TestParseFloat() {
	s.Run("Parse valid float", func() {
		s.Equal(1.234, utils.ParseFloat("1.234"))
	})

	s.Run("Parse invalid float", func() {
		s.Equal(0.0, utils.ParseFloat("abc"))
	})
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(parserSuite))
}
