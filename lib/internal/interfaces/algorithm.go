package interfaces

import "strings"

type Algorithm int

const (
	TokenBucket Algorithm = iota
	FixedWindow
	SlidingWindowLog
	SlidingWindowCounter
	RedisSlidingWindowCounter
)

func ParseAlgorithm(s string) (Algorithm, bool) {
	var algorithmMap = map[string]Algorithm{
		"token-bucket":                 TokenBucket,
		"fixed-window":                 FixedWindow,
		"sliding-window-log":           SlidingWindowLog,
		"sliding-window-counter":       SlidingWindowCounter,
		"redis-sliding-window-counter": RedisSlidingWindowCounter,
	}

	a, ok := algorithmMap[strings.ToLower(s)]

	return a, ok
}

func (d Algorithm) String() string {
	return [...]string{"token-bucket", "fixed-window", "sliding-window-log", "sliding-window-counter", "redis-sliding-window-counter"}[d]
}
