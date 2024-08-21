package serverpool

import (
	"go-load-balancer/constants"
	"go-load-balancer/pkg/lb/internal/server"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type ServerPool struct {
	backends []*server.Server
	current  int64
}

func NewServerPool() *ServerPool {
	obj := &ServerPool{
		backends: make([]*server.Server, 0),
		current:  0,
	}
	go obj.spawnHealthChecker()
	return obj
}

func (s *ServerPool) AddBackend(backend *server.Server) {
	s.backends = append(s.backends, backend)
}

func (s *ServerPool) nextIndex() int64 {
	return (atomic.AddInt64(&s.current, 1) % int64(len(s.backends)))
}

func (s *ServerPool) getNextServer() *server.Server {
	next := s.nextIndex()

	for i := next; i < next+int64(len(s.backends)); i++ {
		idx := i % int64(len(s.backends))

		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreInt64(&s.current, idx)
			}
			return s.backends[idx]
		}
	}
	return nil
}

func (s *ServerPool) ServeNextServer(w http.ResponseWriter, r *http.Request) bool {
	server := s.getNextServer()
	if server != nil {
		server.ServeHTTP(w, r)
		return true
	}
	return false
}

func (s *ServerPool) setAlive(url string, alive bool) {
	for _, b := range s.backends {
		if b.URL() == url {
			b.SetAlive(alive)
		}
	}
}

func (s *ServerPool) healthCheck() {
	log.Default().Println("Starting health check")
	for _, b := range s.backends {
		b.HealthCheck()
	}
	log.Default().Println("Ending health check")
}

func (s *ServerPool) spawnHealthChecker() {
	ticker := time.NewTicker(constants.LoadBalancerTimeout)
	for {
		select {
		case <-ticker.C:
			s.healthCheck()
		}
	}
}
