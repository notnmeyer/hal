package hue

import (
	"os"

	"github.com/amimof/huego"
)

func New() *huego.Bridge {
	var bridge, err = huego.Discover()
	if err != nil {
		panic(err.Error())
	}
	return huego.New(bridge.Host, os.Getenv("HUE_USER"))
}
