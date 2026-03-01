package models 

import (
	"sync"
	"sync/atomic"
	"net/url"
)

type Backend struct {
	URL          *url.URL `json:"url"`
	Alive        bool     `json:"alive"`
	CurrentConns int64    `json:"current_connections"`
	mux          sync.RWMutex
}

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  int64      `json:"current"`
	mux      sync.RWMutex
}

type ProxyConfig struct {
	Port int `json:"port"`
	Strategy string `json:"strategy"`
	HealthCheckFreq	string `json:"health_check_frequency"`
}



func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

 
// ServerPool Methods
func (s *ServerPool) AddBackend(b *Backend) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Backends = append(s.Backends, b)
}


func (s *ServerPool) GetNextValidPeer() *Backend {
	s.mux.RLock()
	defer s.mux.RUnlock()

	count := len(s.Backends)
	if count == 0 {
		return nil
	}

	for i := 0; i < count; i++ {
		next := atomic.AddInt64(&s.Current, 1)
		index := next % int64(count)
		candidate := s.Backends[index]

		if candidate.IsAlive() {
			return candidate
		}
	}
	return nil
}

func (s *ServerPool) GetNextValidPeerLeastConn() *Backend {
	s.mux.RLock()
	defer s.mux.RUnlock()

	var best *Backend
	for _, b := range s.Backends {
		if b.IsAlive() {
			if best == nil || b.GetConns() < best.GetConns() {
				best = b
			}
		}
	}
	return best
}

// SetBackendStatus marks a backend as alive or dead
func (s *ServerPool) SetBackendStatus(url *url.URL, alive bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for _, b := range s.Backends {
		if b.URL.String() == url.String() {
			b.SetAlive(alive)
			break
		}
	}
}

func (s *ServerPool) GetBackends() []*Backend {
	s.mux.RLock()
	defer s.mux.RUnlock()
	
	list := make([]*Backend, len(s.Backends))
	copy(list, s.Backends)
	return list
}

// RemoveBackend removes a backend from the pool
func (s *ServerPool) RemoveBackend(u *url.URL) {
	s.mux.Lock()
	defer s.mux.Unlock()

	for i, b := range s.Backends {
		if b.URL.String() == u.String() {
			s.Backends = append(s.Backends[:i], s.Backends[i+1:]...)
			return
		}
	}
}

func (b *Backend) GetConns() int64 {
	return atomic.LoadInt64(&b.CurrentConns)
}

func (b *Backend) IncConns() {
	atomic.AddInt64(&b.CurrentConns, 1)
}

func (b *Backend) DecConns() {
	atomic.AddInt64(&b.CurrentConns, -1)
}
