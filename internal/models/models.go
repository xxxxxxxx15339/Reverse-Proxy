package models 

	"time"
	"sync"
	"sync/atomic"
	"net/url"
)

type Backend struct {
	URL *url.URL 	`json:"url"`
	Alive bool 		`json:"alive"`
	CurrentConns 	int64 	`json:"current_connections"`
	mux sync.RWMutex 
}

type ServerPool struct {
	Backends []*Backend  `json:"backends"`
	Current int64 `json:"current"`
	mux sync.RWMutex
}	

type ProxyConfig struct {
	Port int `json:"port"`
	Strategy string `json:"strategy"`
	HealthCheckFreq	string `json:"health_check_frequency"`
}



// Thread-Safe Methods 
// for SetAlive, we are changing the data so we need to use a Write Lock (Exclusive)
func (b *Backend) SetAlive(alive bool) {
	// Lock the Backend
	b.mux.Lock()
	defer b.mux.Unlock()
	// Set the alive status 
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


// With Atomic Increment, The CPU treats the read-modify-write operation as one single, uninterruptible step.
func (s *ServerPool) GetNextValidPeer() *Backend {
	s.mux.RLock()
	defer s.mux.RUnlock()

	// Loop through the backends to find an available one
	// We use the length of the slice to determine the cycle
	count := len(s.Backends)
	if count == 0 {
		return nil
	}

	// Try to find a valid peer, cycling through the list if necessary
	// We loop 'count' times to ensure we check everyone once
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

// GetBackends returns a slice of backends safely
func (s *ServerPool) GetBackends() []*Backend {
	s.mux.RLock()
	defer s.mux.RUnlock()
	
	// Create a new slice to avoid race conditions on the slice header
	// The elements are pointers, which is fine since Backend methods are thread-safe
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
