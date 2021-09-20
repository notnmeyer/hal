package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/notnmeyer/hal/internal/hue"
	l "github.com/notnmeyer/hal/internal/logger"
	"github.com/notnmeyer/hal/internal/plex"
)

func main() {
	logger := l.New()
	logger.Info("Logger initialized")

	bridge := hue.New()
	logger.Info("Hue client initialized")

	// routes
	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/healthcheck", healthHandler)
	http.HandleFunc("/plex", plex.WebhookHandler(logger, bridge))

	// start the server
	listenOn := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	logger.Infof("Listening on %s", listenOn)
	http.ListenAndServe(listenOn, http.DefaultServeMux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
