package plex

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/amimof/huego"
	"github.com/hekmon/plexwebhooks"
	"go.uber.org/zap"

	"github.com/notnmeyer/hal/internal/config"
	"github.com/notnmeyer/hal/internal/hue"
	"github.com/notnmeyer/hal/internal/redisClient"
)

var (
	cfg   = *make(config.Config).New()
	redis = *redisClient.New()
)

func WebhookHandler(logger *zap.SugaredLogger, bridge *huego.Bridge) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// light.control is a bool indicating whether to react to webhook
		// right now.
		enabled, err := strconv.ParseBool(
			redis.HGet("light.control", "enabled").Val(),
		)

		if err != nil {
			logger.Errorf("Failed to parse enabled value: %s", err)
			return
		}

		if !enabled {
			logger.Info("Light control is disabled. Skipping event...")
			return
		}

		// see docs/plex_payload_example.md
		payload, err := getRequestPayload(w, r)
		if err != nil {
			logger.Errorf("Error getting request payload: %s", err.Error())
			return
		}

		if payload.Player.Title == cfg["BONUS_ROOM_PLEX_CLIENT_NAME"] {
			err := hue.EventHandler(logger, bridge, payload, cfg["BONUS_ROOM_HUE"])
			if err != nil {
				logger.Errorf("Error handling event: %s", err.Error())
			}
		}
	}
}

func getRequestPayload(w http.ResponseWriter, r *http.Request) (*plexwebhooks.Payload, error) {
	multiPartReader, err := r.MultipartReader()
	if err != nil {
		if err == http.ErrNotMultipart || err == http.ErrMissingBoundary {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, wErr := w.Write([]byte(err.Error()))
		if wErr != nil {
			err = fmt.Errorf("request error: %v | write error: %v", err, wErr)
		}
		return nil, fmt.Errorf("can't create a multipart reader from request: %s", err)
	}

	payload, _, err := plexwebhooks.Extract(multiPartReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, wErr := w.Write([]byte(err.Error()))
		if wErr != nil {
			err = fmt.Errorf("request error: %v | write error: %v", err, wErr)
		}
		return nil, fmt.Errorf("can't create a multipart reader from request: %s", err)
	}
	return payload, nil
}
