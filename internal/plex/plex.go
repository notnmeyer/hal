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

		// Create the multi part reader
		multiPartReader, err := r.MultipartReader()
		if err != nil {
			// Detect error type for the http answer
			if err == http.ErrNotMultipart || err == http.ErrMissingBoundary {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			// Try to write the error as http body
			_, wErr := w.Write([]byte(err.Error()))
			if wErr != nil {
				err = fmt.Errorf("request error: %v | write error: %v", err, wErr)
			}
			// Log the error
			logger.Info("can't create a multipart reader from request:", err)
			return
		}
		// Use the multipart reader to parse the request body
		payload, _, err := plexwebhooks.Extract(multiPartReader)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// Try to write the error as http body
			_, wErr := w.Write([]byte(err.Error()))
			if wErr != nil {
				err = fmt.Errorf("request error: %v | write error: %v", err, wErr)
			}
			// Log the error
			logger.Info("can't create a multipart reader from request:", err)
			return
		}
		// spew.Dump("%#v\n", *payload)

		if payload.Player.Title == os.Getenv("PLEX_CLIENT_NAME") {
			group := hue.InitBonusRoom(logger, bridge)
			switch event := payload.Event; event {
			case "media.play", "media.resume":
				hue.BonusRoomOn(logger, group)
			case "media.stop", "media.pause":
				hue.BonusRoomOff(logger, group)
			}
		}
	}
}
