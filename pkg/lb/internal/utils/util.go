package utils

import (
	"go-load-balancer/constants"
	"net/http"
)

func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(constants.Attempts).(int); ok {
		return attempts
	}
	return 0
}

func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(constants.Retry).(int); ok {
		return retry
	}
	return 0
}
