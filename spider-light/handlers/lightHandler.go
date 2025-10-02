package handlers

import (
	"net/http"
	"spider-light/globals"
	"spider-light/hardwares/ledactions"
)

func LightHandler(w http.ResponseWriter, r *http.Request) {
	lightType := r.URL.Query().Get("type")
	switch {
	case lightType == "none":
		globals.LedControl.SetAction(ledactions.NoneLight(globals.LedControl))

	case lightType == "police_calm":
		globals.LedControl.SetAction(ledactions.PoliceCalmLight(globals.LedControl))

	case lightType == "police_warn":
		globals.LedControl.SetAction(ledactions.PoliceWarmlight(globals.LedControl))

	case lightType == "warning":
		globals.LedControl.SetAction(ledactions.WarningLight(globals.LedControl))

	default:
		globals.LedControl.SetAction(ledactions.NoneLight(globals.LedControl))
		return
	}
	w.WriteHeader(http.StatusOK)
}
