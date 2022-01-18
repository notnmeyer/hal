package main

import (
	"fmt"
	"net/http"

	"github.com/amimof/huego"
	"github.com/notnmeyer/hal/internal/config"
	"github.com/notnmeyer/hal/internal/hue"
	l "github.com/notnmeyer/hal/internal/logger"
	"github.com/notnmeyer/hal/internal/plex"
	"go.uber.org/zap"
)

var cfg config.Config
var logger *zap.SugaredLogger
var bridge *huego.Bridge

func init() {
	logger = l.New()
	logger.Info("Logger initialized")

	cfg = *cfg.New()
	logger.Info("Config initialized")

	bridge = hue.New()
	logger.Info("Hue client initialized")
}

func main() {
	// routes
	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/healthcheck", healthHandler)

	http.HandleFunc("/plex", plex.WebhookHandler(logger, bridge))
	http.HandleFunc("/plex/configure", plex.ConfigureHandler)

	// start the server
	listenOn := fmt.Sprintf(":%s", cfg["PORT"])
	logger.Infof("Listening on %s", listenOn)
	http.ListenAndServe(listenOn, http.DefaultServeMux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
