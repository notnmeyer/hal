package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	cfg := *Config{}.New()

	if _, ok := cfg["BONUS_ROOM_HUE"]; !ok {
		t.Errorf("BONUS_ROOM_HUE env var is not set")
	}

	if _, ok := cfg["BONUS_ROOM_PLEX_CLIENT_NAME"]; !ok {
		t.Errorf("BONUS_ROOM_PLEX_CLIENT_NAME env var is not set")
	}

	if _, ok := cfg["HUE_USER"]; !ok {
		t.Errorf("HUE_USER env var is not set")
	}

	if _, ok := cfg["PORT"]; !ok {
		t.Errorf("PORT env var is not set")
	}

	if _, ok := cfg["REDIS_URL"]; !ok {
		t.Errorf("REDIS_URL env var is not set")
	}
}
