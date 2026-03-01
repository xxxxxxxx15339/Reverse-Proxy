package proxy

import (
	"net/http"
	"net/http/httputil"
	"reverse-proxy/internal/models"
	"reverse-proxy/internal/server"
)

type ProxyHandler struct {
	Pool     *models.ServerPool
	Balancer server.LoadBalancer
	Strategy string
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var peer *models.Backend
	if h.Strategy == "least-conn" {
		peer = h.Pool.GetNextValidPeerLeastConn()
	} else {
		peer = h.Pool.GetNextValidPeer()
	}

	if peer == nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(peer.URL)
	peer.IncConns()
	defer peer.DecConns()

	proxy.ServeHTTP(w, r)
}
