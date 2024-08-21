package server

import (
	"go-load-balancer/constants"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	url          *url.URL
	alive        bool
	mutex        *sync.RWMutex
	reverseProxy *httputil.ReverseProxy
}

func NewServer(url *url.URL, reverseProxy *httputil.ReverseProxy) *Server {
	return &Server{
		url:          url,
		alive:        true,
		mutex:        &sync.RWMutex{},
		reverseProxy: reverseProxy,
	}
}

func (s *Server) IsAlive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.alive
}

func (s *Server) SetAlive(alive bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.alive = alive
}

func (s *Server) URL() string {
	return s.url.String()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.reverseProxy.ServeHTTP(w, r)
}

func (s *Server) HealthCheck() bool {
	conn, err := net.DialTimeout("tcp", s.url.Host, constants.ServerHealthTimeout)
	if err != nil {
		s.SetAlive(false)
		return false
	}

	conn.Close()
	s.SetAlive(true)
	return true
}
