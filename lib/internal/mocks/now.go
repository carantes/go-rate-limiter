package mocks

import "time"

// Now is a mockable version of time.Now, it is used to mock the current time in tests
var Now = time.Now
