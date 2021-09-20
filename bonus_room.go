package main

import (
	"fmt"

	"github.com/amimof/huego"
	"go.uber.org/zap"
)

func initBonusRoom(logger *zap.SugaredLogger, bridge *huego.Bridge) huego.Group {
	var bonusRoomGroup huego.Group
	var bonusRoomGroupName = "Bonus Room"

	groups, err := bridge.GetGroups()
	if err != nil {
		fmt.Println(err.Error())
	}

	// find the specific group
	for _, group := range groups {
		if group.Name == bonusRoomGroupName {
			bonusRoomGroup = group
			logger.Debugf("Group `%s` found", group.Name)
		}
	}

	return bonusRoomGroup
}

func bonusRoomOn(logger *zap.SugaredLogger, group huego.Group) {
	logger.Infof("Turning off lights in Bonus Room: %s", group.Name)
	group.Off()
}

func bonusRoomOff(logger *zap.SugaredLogger, group huego.Group) {
	logger.Infof("Turning on lights in Bonus Room: %s", group.Name)
	group.Bri(10)
	group.On()
}
