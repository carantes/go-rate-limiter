package utils

import "strconv"

// ParseInt parse string to int
func ParseInt(s string) int {
	r, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return r
}

// ParseFloat parse string to float64
func ParseFloat(s string) float64 {
	r, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0
	}

	return r
}
