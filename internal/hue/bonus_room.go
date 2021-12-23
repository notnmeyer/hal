package hue

import (
	"fmt"
	"os"

	"github.com/amimof/huego"
	"go.uber.org/zap"
)

func InitBonusRoom(logger *zap.SugaredLogger, bridge *huego.Bridge) huego.Group {
	var bonusRoomGroup huego.Group

	groups, err := bridge.GetGroups()
	if err != nil {
		fmt.Println(err.Error())
	}

	// find the specific group
	for _, group := range groups {
		if group.Name == os.Getenv("HUE_MEDIA_ROOM_GROUP") {
			bonusRoomGroup = group
			break
		}
	}

	// TODO: do something with failure case when no matching group
	return bonusRoomGroup
}

func BonusRoomOn(logger *zap.SugaredLogger, group huego.Group) {
	logger.Infof("Turning off lights in Bonus Room: %s", group.Name)
	group.Off()
}

func BonusRoomOff(logger *zap.SugaredLogger, group huego.Group) {
	logger.Infof("Turning on lights in Bonus Room: %s", group.Name)
	group.Bri(60)
	group.On()
}
