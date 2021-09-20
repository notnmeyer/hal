package main

import (
	"net/http"

	"github.com/notnmeyer/hal/internal/hue"
	l "github.com/notnmeyer/hal/internal/logger"
	"github.com/notnmeyer/hal/internal/plex"
)

func main() {
	logger := l.New()
	logger.Info("Logger initialized")

	bridge := hue.New()
	logger.Info("Hue client initialized")

	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/healthcheck", healthHandler)
	http.HandleFunc("/plex", plex.WebhookHandler(logger, bridge))
	http.ListenAndServe(":7095", http.DefaultServeMux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
