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
