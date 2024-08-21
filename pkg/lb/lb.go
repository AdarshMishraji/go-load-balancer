package lb

import (
	"context"
	"go-load-balancer/constants"
	"go-load-balancer/pkg/lb/internal/server"
	"go-load-balancer/pkg/lb/internal/serverpool"
	"go-load-balancer/pkg/lb/internal/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var serverPool = serverpool.NewServerPool()

func BalanceLoad(w http.ResponseWriter, r *http.Request) {
	attempts := utils.GetAttemptsFromContext(r)

	if attempts > constants.MaxAttempts {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	isServed := serverPool.ServeNextServer(w, r)
	if !isServed {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}
}

func RegisterServerWithReverseProxy(url *url.URL) {
	proxy := httputil.NewSingleHostReverseProxy(url)
	server := server.NewServer(url, proxy)

	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		retries := utils.GetRetryFromContext(request)
		if retries > constants.MaxServerRetries {
			server.SetAlive(false)

			attempts := utils.GetAttemptsFromContext(request)
			ctx := context.WithValue(request.Context(), constants.Attempts, attempts+1)
			BalanceLoad(writer, request.WithContext(ctx))

			http.Error(writer, "service not available", http.StatusServiceUnavailable)
			return
		}

		ticker := time.NewTicker(constants.RetryTimeout)
		select {
		case <-ticker.C:
			ctx := context.WithValue(request.Context(), constants.Retry, retries+1)
			server.ServeHTTP(writer, request.WithContext(ctx))
		}
	}

	serverPool.AddBackend(server)
}
