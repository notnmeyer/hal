package plex

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/notnmeyer/hal/internal/redisClient"
)

func ConfigureHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	_, exists := query["enabled"]

	if !exists {
		msg := "`enabled` query parameter is missing"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	enabled, err := strconv.ParseBool(query["enabled"][0])
	if err != nil {
		msg := "`enabled` query parameter is not a boolean"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	redisClient.New().HSet("light.control", "enabled", enabled)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%t", enabled)))
}
