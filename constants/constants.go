package constants

import "time"

const (
	ServerHealthTimeout = 2 * time.Second
	LoadBalancerTimeout = time.Minute
	RetryTimeout        = 10 * time.Millisecond
	MaxAttempts         = 3
	MaxServerRetries    = 3
	Attempts            = "attempts"
	Retry               = "retry"
)
