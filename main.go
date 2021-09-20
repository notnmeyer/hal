package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/amimof/huego"
	"github.com/hekmon/plexwebhooks"
	"go.uber.org/zap"
)

func initHue() *huego.Bridge {
	var bridge, err = huego.Discover()
	if err != nil {
		panic(err.Error())
	}
	return huego.New(bridge.Host, os.Getenv("HUE_USER"))
}

func initLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err.Error())
	}
	sugar := logger.Sugar()
	return sugar
}

func main() {
	logger := initLogger()
	logger.Info("Logger initialized")

	bridge := initHue()
	logger.Info("Hue client initialized")

	http.HandleFunc("/", webhookHandler(logger, bridge))
	http.ListenAndServe(":7095", http.DefaultServeMux)
}

func webhookHandler(logger *zap.SugaredLogger, bridge *huego.Bridge) http.HandlerFunc {
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

		//
		// Bonus Room
		//
		if payload.Player.Title == "Bonus room" { // "Bonus room" is correct for the Plex player
			// get the hue group
			group := initBonusRoom(logger, bridge)
			switch event := payload.Event; event {
			case "media.play", "media.resume":
				bonusRoomOn(logger, group)
			case "media.stop", "media.pause":
				bonusRoomOff(logger, group)
			}
		}
	}
}
