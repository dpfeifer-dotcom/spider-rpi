package handlers

import (
	"net/http"
	"spider-light/globals"
	"spider-light/hardwares"
)

func BulbHandler(w http.ResponseWriter, r *http.Request) {
	hardwares.SwitchBulbLight(globals.Bulb)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
