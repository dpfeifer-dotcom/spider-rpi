package main

import (
	"fmt"
	"log"
	"net/http"
	"spider-light/globals"
	"spider-light/handlers"
	"spider-light/hardwares"
	"spider-light/hardwares/ledactions"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func main() {
	if _, err := host.Init(); err != nil {
		fmt.Println("Hiba az inicializálás során:", err)
		return
	}

	globals.Bulb = gpioreg.ByName("GPIO5")
	if globals.Bulb == nil {
		fmt.Println("A GPIO pin nem található!")
		return
	}
	globals.Bulb.Out(gpio.Low) // vagy gpio.High, ahogy akarod az alaphelyzetet

	globals.LedControl = hardwares.NewLedControll()
	globals.LedControl.LedService = ledactions.NoneLight(globals.LedControl)

	http.HandleFunc("/light", handlers.LightHandler)
	http.HandleFunc("/bulb", handlers.BulbHandler)

	log.Println("server listening on :8080")
	go http.ListenAndServe("0.0.0.0:8080", nil)

	select {}
}
