package plex

import (
	"fmt"
	"net/http"
	"os"

	"github.com/amimof/huego"
	"github.com/hekmon/plexwebhooks"
	"go.uber.org/zap"

	"github.com/notnmeyer/hal/internal/hue"
)

func WebhookHandler(logger *zap.SugaredLogger, bridge *huego.Bridge) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// see docs/plex_payload_example.md
		payload, err := getRequestPayload(w, r)
		if err != nil {
			logger.Errorf("Error getting request payload: %s", err.Error())
			return
		}

		if payload.Player.Title == os.Getenv("BONUS_ROOM_PLEX_CLIENT_NAME") {
			err := hue.EventHandler(logger, bridge, payload, os.Getenv("BONUS_ROOM_HUE"))
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
