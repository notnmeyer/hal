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

		payload, err := getRequestPayload(logger, w, r)
		if err != nil {
			logger.Errorf("Error getting request payload: %s", err.Error())
			return
		}

		// spew.Dump("%#v\n", *payload)

		if payload.Player.Title == os.Getenv("PLEX_CLIENT_NAME") {
			group := hue.InitBonusRoom(logger, bridge)
			switch event := payload.Event; event {
			case "media.play", "media.resume":
				hue.BonusRoomOff(logger, group)
				logger.Infof("handling `%s` event", event)
			case "media.stop", "media.pause":
				hue.BonusRoomOn(logger, group)
				logger.Infof("handling `%s` event", event)
			default:
				logger.Infof("ignoring `%s` event", event)
			}
		}
	}
}

func getRequestPayload(logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) (*plexwebhooks.Payload, error) {
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
