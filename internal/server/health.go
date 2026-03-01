package server

import (
	"log"
	"net"
	"reverse-proxy/internal/models"
	"time"
)

func Ping(b *models.Backend) bool {
	conn, err := net.DialTimeout("tcp", b.URL.Host, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func RunHealthCheck(s *models.ServerPool, freq time.Duration) {
	ticker := time.NewTicker(freq)
	for range ticker.C {
		for _, b := range s.GetBackends() {
			alive := Ping(b)
			b.SetAlive(alive)
			if !alive {
				log.Printf("Backend %s is DOWN", b.URL.String())
			}
		}
	}
}
