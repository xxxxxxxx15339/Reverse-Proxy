package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reverse-proxy/internal/models"
)

type AdminAPI struct {
	Pool *models.ServerPool
}

func (a *AdminAPI) GetStatus(w http.ResponseWriter, r *http.Request) {
	backends := a.Pool.GetBackends()
	var active int
	for _, b := range backends {
		if b.IsAlive() {
			active++
		}
	}

	res := map[string]interface{}{
		"total_backends":  len(backends),
		"active_backends": active,
		"backends":        backends,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (a *AdminAPI) AddBackend(w http.ResponseWriter, r *http.Request) {
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := url.Parse(body.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.Pool.AddBackend(&models.Backend{
		URL:   u,
		Alive: true,
	})
	w.WriteHeader(http.StatusCreated)
}

func (a *AdminAPI) RemoveBackend(w http.ResponseWriter, r *http.Request) {
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := url.Parse(body.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.Pool.RemoveBackend(u)
	w.WriteHeader(http.StatusNoContent)
}
