package interfaces

// Custom Error
type RateLimitError struct {
	Message string
}

func (r *RateLimitError) Error() string {
	panic(r.Message)
}
