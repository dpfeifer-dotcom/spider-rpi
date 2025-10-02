package main

import (
	"fmt"
	"log"
	"net/http"
	"spider-backend/globals"
	"spider-backend/handlers"
	"spider-backend/hardwares"
	"spider-backend/hardwares/ledactions"
	cache "spider-backend/memorycache"
	"spider-backend/services"
	"time"

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

	globals.CPUUsageCache = cache.NewMemoryListCache("cpu_usage-storage").SetLimits(100).SetDefaultData(0.0).StartWorkers().Verbose()
	globals.CPUTempCache = cache.NewMemoryListCache("cpu_temp-storage").SetDefaultData(0.0).StartWorkers().Verbose()
	globals.MemoryUsageCache = cache.NewMemoryListCache("mem_usage-storage").SetDefaultData(0.0).StartWorkers().Verbose()

	globals.SystemCache = cache.NewMemoryCache("system").StartWorkers().Verbose()
	globals.LedControl = hardwares.NewLedControll()
	globals.LedControl.LedService = ledactions.NoneLight(globals.LedControl)
	globals.Controller = hardwares.NewController(0)

	//services.SystemCheckService(globals.SystemCache, "cpu_usage", hardwares.GetCpuUsage, 1)
	//services.SystemCheckService(globals.SystemCache, "cpu_temp", hardwares.GetCpuTemperature, 1)
	//services.SystemCheckService(globals.SystemCache, "mem_usage", hardwares.GetMemoryUsage, 1)

	services.SystemListCheckService(globals.CPUUsageCache, hardwares.GetCpuUsage, 1)
	services.SystemListCheckService(globals.CPUTempCache, hardwares.GetCpuTemperature, 1)
	services.SystemListCheckService(globals.MemoryUsageCache, hardwares.GetMemoryUsage, 1)

	/* go func() {
		for {
			time.Sleep(3 * time.Second)

			fmt.Printf("CPU kihasználtság: %.2f%%\n", globals.SystemCache.Get("cpu_usage").Value.(float64))
			fmt.Printf("CPU hőmérséklet: %.2f°C\n", globals.SystemCache.Get("cpu_temp").Value.(float64))
			fmt.Printf("Memória használat: %.2f%%\n\n", globals.SystemCache.Get("mem_usage").Value.(float64))
			fmt.Println(globals.Controller.State())
		}
	}() */
	btnAwasPressed := false
	btnBwasPressed := false
	btnXwasPressed := false
	btnYwasPressed := false
	btnLBwasPressed := false
	btnRsBwasPressed := false
	btnStartwasPressed := false
	btnBackwasPressed := false
	btnGuidewasPressed := false
	btnLSBwasPressed := false
	btnRSBwasPressed := false

	btnDPADRightwasPressed := false
	btnDPADLeftwasPressed := false
	btnDPADUpasPressed := false
	btnDPADDownwasPressed := false

	http.HandleFunc("/light", handlers.LightHandler)
	http.HandleFunc("/bulb", handlers.BulbHandler)
	http.HandleFunc("/cpu_usage", handlers.CPUSUsageSensorHandler)
	http.HandleFunc("/cpu_usage_all", handlers.CPUSUsageAllSensorHandler)
	http.HandleFunc("/cpu_temp", handlers.CPUSTempSensorHandler)
	http.HandleFunc("/cpu_temp_all", handlers.CPUSTempAllSensorHandler)

	http.HandleFunc("/mem_usage", handlers.MemoryUsageSensorHandler)
	http.HandleFunc("/mem_usage_all", handlers.MemoryUsageAllSensorHandler)

	log.Println("server listening on :8080")
	go http.ListenAndServe("0.0.0.0:8080", nil)

	go func() {
		for {
			state, err := globals.Controller.State()
			if err != nil {
				log.Printf("Hiba a vezérlő állapotának lekérdezésekor: %v\n", err)
				continue
			}
			hardwares.ButtonPressed(state.A(), &btnAwasPressed, func() { globals.LedControl.SetAction(ledactions.NoneLight(globals.LedControl)) })
			hardwares.ButtonPressed(state.B(), &btnBwasPressed, func() { globals.LedControl.SetAction(ledactions.PoliceWarmlight(globals.LedControl)) })
			hardwares.ButtonPressed(state.X(), &btnXwasPressed, func() { globals.LedControl.SetAction(ledactions.WarningLight(globals.LedControl)) })
			hardwares.ButtonPressed(state.Y(), &btnYwasPressed, func() { globals.LedControl.SetAction(ledactions.PoliceCalmLight(globals.LedControl)) })

			hardwares.ButtonPressed(state.LB(), &btnLBwasPressed, func() { log.Println("LB pressed") })
			hardwares.ButtonPressed(state.RB(), &btnRsBwasPressed, func() { hardwares.SwitchBulbLight(globals.Bulb) })

			hardwares.ButtonPressed(state.Start(), &btnStartwasPressed, func() { log.Println("Start pressed") })
			hardwares.ButtonPressed(state.Back(), &btnBackwasPressed, func() { log.Println("Back pressed") })
			hardwares.ButtonPressed(state.Guide(), &btnGuidewasPressed, func() { log.Println("Guide pressed") })

			hardwares.ButtonPressed(state.LSB(), &btnLSBwasPressed, func() { log.Println("LSB pressed") })
			hardwares.ButtonPressed(state.RSB(), &btnRSBwasPressed, func() { log.Println("RSB pressed") })

			hardwares.ButtonPressed(state.DPadDown(), &btnDPADDownwasPressed, func() { log.Println("DPAD Down pressed") })
			hardwares.ButtonPressed(state.DPadUp(), &btnDPADUpasPressed, func() { log.Println("DPAD Up pressed") })
			hardwares.ButtonPressed(state.DPadRight(), &btnDPADRightwasPressed, func() { log.Println("DPAD Right pressed") })
			hardwares.ButtonPressed(state.DPadLeft(), &btnDPADLeftwasPressed, func() { log.Println("DPAD Left pressed") })

			time.Sleep(time.Second / 30)
		}

	}()
	select {}
}
