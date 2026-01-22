package models 

import (
	"time"
	"sync"
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
}	

type ProxyConfig struct {
	Port int `json:"port"`
	Strategy string `json:"strategy"`
	HealthCheckFreq	time.Duration `json:"health_check_frequency"`
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
func (s *ServerPool) Addbackend(b *Backend) {
	s.Backends.append(s.Backends, b)
}


// With Atomic Increment, The CPU treats the read-modify-write operation as one single, uninterruptible step.
func (s *ServerPool) GetNextValidPeer() *Backend {
	for i := 0; i < len(s.Backends); i++ {	
			next := atomic.AddInt64(&s.Current, 1) // This actually updates the Current in the ServerPool, directly by taking a pointer to s.Current and incrementing it.
			index := next % int64(len(s.Backends))

			candidate := s.Backends[index]
				
				if candidate.IsAlive() {
					return candidate
				}

	}
	return nil
}
