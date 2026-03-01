package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reverse-proxy/internal/api"
	"reverse-proxy/internal/models"
	"reverse-proxy/internal/proxy"
	"reverse-proxy/internal/server"
	"time"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var config models.ProxyConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	pool := &models.ServerPool{}
	
	freq, err := time.ParseDuration(config.HealthCheckFreq)
	if err != nil {
		freq = 10 * time.Second
	}

	go server.RunHealthCheck(pool, freq)

	proxyHandler := &proxy.ProxyHandler{
		Pool:     pool,
		Strategy: config.Strategy,
	}

	adminAPI := &api.AdminAPI{Pool: pool}
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/status", adminAPI.GetStatus)
	adminMux.HandleFunc("/backends", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			adminAPI.AddBackend(w, r)
		case http.MethodDelete:
			adminAPI.RemoveBackend(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	go func() {
		log.Printf("Admin API starting on :8081")
		if err := http.ListenAndServe(":8081", adminMux); err != nil {
			log.Fatal(err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", config.Port)
	}

	log.Printf("Proxy server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, proxyHandler); err != nil {
		log.Fatal(err)
	}
}
