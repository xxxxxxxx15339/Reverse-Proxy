package server

import (
	"reverse-proxy/internal/models"
	"net/url"
)

type LoadBalancer interface {
	GetNextValidPeer() *models.Backend
	AddBackend(backend *models.Backend)
	SetBackendStatus(uri *url.URL, alive bool)
}
