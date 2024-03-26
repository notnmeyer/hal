package hue

import (
	"image/color"
	"net/http"
	"time"

	"github.com/amimof/huego"
	"go.uber.org/zap"
)

var (
	colors = map[string][]color.Color{
		"default": {
			color.RGBA{255, 0, 0, 255},
			color.RGBA{0, 255, 0, 255},
			color.RGBA{0, 0, 255, 255},
		},
		"space": {
			color.RGBA{0, 0, 255, 255},
			color.RGBA{255, 255, 255, 255},
		},
	}
)

func Party(logger *zap.SugaredLogger, bridge *huego.Bridge, groupName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger.Infof("party time!")

		group, err := GetGroup(bridge, groupName)
		if err != nil {
			logger.Errorf("error getting group: %s", err)
			return
		}

		// this is definitely not the right way to do this. never stops, and the request never completes
		for {
			for _, color := range colors["space"] {
				xy, _ := huego.ConvertRGBToXy(color)
				logger.Infof("setting color to: %v", xy)
				group.Xy(xy)
				time.Sleep(1000 * time.Millisecond)
			}
		}
	}
}
