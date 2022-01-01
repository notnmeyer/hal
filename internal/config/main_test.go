package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := *Config{}.New()

	if cfg["BONUS_ROOM_HUE"] != os.Getenv("BONUS_ROOM_HUE") {
		t.Errorf("BONUS_ROOM_HUE env var is not set")
	}
}
