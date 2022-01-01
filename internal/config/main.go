package config

import "os"

type Config map[string]string

func (cfg Config) New() *Config {
	cfg["BONUS_ROOM_HUE"] = os.Getenv("BONUS_ROOM_HUE")
	cfg["BONUS_ROOM_PLEX_CLIENT_NAME"] = os.Getenv("BONUS_ROOM_PLEX_CLIENT_NAME")
	cfg["HUE_USER"] = os.Getenv("HUE_USER")
	cfg["PORT"] = os.Getenv("PORT")
	cfg["REDIS_URL"] = os.Getenv("REDIS_URL")
	return &cfg
}
