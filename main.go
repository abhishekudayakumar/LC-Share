package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"time"
)

const (
	appName    = "LC Share"
	appVersion = "0.2.0"
	listenAddr = "127.0.0.1:8765"
)

//go:embed web/*
var webFiles embed.FS

type apiResponse struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

func writeJSON(w http.ResponseWriter, statusCode int, payload apiResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Service: "lc-share",
		Version: appVersion,
		Status:  "ok",
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Service: "lc-share",
		Version: appVersion,
		Status:  "online",
	})
}

func newHandler() (http.Handler, error) {
	web, err := fs.Sub(webFiles, "web")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/status", statusHandler)
	mux.Handle("/", http.FileServer(http.FS(web)))

	return mux, nil
}

func main() {
	handler, err := newHandler()
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("%s v%s", appName, appVersion)
	log.Printf("Dashboard: http://%s", listenAddr)
	log.Printf("Health API: http://%s/api/health", listenAddr)
	log.Printf("Status API: http://%s/api/status", listenAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
