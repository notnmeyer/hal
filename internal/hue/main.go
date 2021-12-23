package hue

import (
	"fmt"
	"os"

	"github.com/amimof/huego"
	"github.com/hekmon/plexwebhooks"
	"go.uber.org/zap"
)

func New() *huego.Bridge {
	var bridge, err = huego.Discover()
	if err != nil {
		panic(err.Error())
	}
	return huego.New(bridge.Host, os.Getenv("HUE_USER"))
}

func GetGroup(bridge *huego.Bridge, groupName string) (huego.Group, error) {
	groups, err := bridge.GetGroups()
	if err != nil {
		return huego.Group{}, fmt.Errorf("error getting groups: %s", err.Error())
	}

	// find the specific group
	for _, group := range groups {
		if group.Name == groupName {
			return group, nil
		}
	}

	return huego.Group{}, fmt.Errorf("no Hue group found with name: %s", groupName)
}

func EventHandler(
	logger *zap.SugaredLogger,
	bridge *huego.Bridge,
	payload *plexwebhooks.Payload,
	groupName string,
) error {
	group, err := GetGroup(bridge, groupName)
	if err != nil {
		return err
	}

	switch event := payload.Event; event {
	case "media.play", "media.resume":
		logger.Infof("handling `%s` event", event)
		GroupOff(&group)
	case "media.stop", "media.pause":
		logger.Infof("handling `%s` event", event)
		GroupOn(&group, 60)
	default:
		logger.Infof("ignoring `%s` event", event)
	}

	return nil
}

func GroupOn(group *huego.Group, brightness uint8) {
	group.Bri(brightness)
	group.On()
}

func GroupOff(group *huego.Group) {
	group.Off()
}

// if a video is paused before sleeping the apple tv a stop event will be sent a few minutes later.
// if everything is turned off (including lights) the stop event is handled and the lights
// are mysteriously turned on.
