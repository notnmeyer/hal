package hue

import (
	"context"
	"fmt"
	"os"

	"github.com/amimof/huego"
	"github.com/go-redis/redis/v8"
	"github.com/hekmon/plexwebhooks"
	"go.uber.org/zap"
)

var ctx = context.Background()

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
		GroupOff(&group, string(event))
	case "media.stop":
		// if the previous event was "media.pause", then the stop event is
		// may be received as a result of the device turning off. in any case,
		// we don't need to turn the lights back on in that situation.
		if readHistory(group.Name) == "media.pause" {
			logger.Infof("skipping `%s` event because previous event was `media.pause`", event)
			return nil
		}
		fallthrough
	case "media.pause":
		logger.Infof("handling `%s` event", event)
		GroupOn(&group, 60, string(event))
	default:
		logger.Infof("ignoring `%s` event", event)
	}

	return nil
}

func GroupOn(group *huego.Group, brightness uint8, event string) {
	group.Bri(brightness)
	group.On()
	updateHistory(group.Name, event)
}

func GroupOff(group *huego.Group, event string) {
	group.Off()
	updateHistory(group.Name, event)
}

func redisClient() *redis.Client {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err.Error())
	}
	return redis.NewClient(opts)
}

// writes the last plex event type to redis for the hue group
func updateHistory(hueGroupName, plexEvent string) {
	rdb := redisClient()
	rdb.HSet(ctx, "history", hueGroupName, plexEvent)
}

func readHistory(hueGroupName string) string {
	rdb := redisClient()
	return rdb.HGet(ctx, "history", hueGroupName).Val()
}
