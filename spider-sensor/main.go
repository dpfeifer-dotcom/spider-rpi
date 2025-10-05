package main

import (
	"log"
	"net/http"
	"spider-sensor/globals"
	"spider-sensor/handlers"
	"spider-sensor/hardwares"
	cache "spider-sensor/memorycache"
	"spider-sensor/services"
)

func main() {
	globals.CPUUsageCache = cache.NewMemoryListCache("cpu_usage-storage").SetLimits(100).SetDefaultData(0.0).StartWorkers()
	globals.CPUTempCache = cache.NewMemoryListCache("cpu_temp-storage").SetDefaultData(0.0).StartWorkers()
	globals.MemoryUsageCache = cache.NewMemoryListCache("mem_usage-storage").SetDefaultData(0.0).StartWorkers()

	services.SystemListCheckService(globals.CPUUsageCache, hardwares.GetCpuUsage, 1)
	services.SystemListCheckService(globals.CPUTempCache, hardwares.GetCpuTemperature, 1)
	services.SystemListCheckService(globals.MemoryUsageCache, hardwares.GetMemoryUsage, 1)

	http.HandleFunc("/cpu_usage", handlers.CPUSUsageSensorHandler)
	http.HandleFunc("/cpu_usage_all", handlers.CPUSUsageAllSensorHandler)
	http.HandleFunc("/cpu_temp", handlers.CPUSTempSensorHandler)
	http.HandleFunc("/cpu_temp_all", handlers.CPUSTempAllSensorHandler)

	http.HandleFunc("/mem_usage", handlers.MemoryUsageSensorHandler)
	http.HandleFunc("/mem_usage_all", handlers.MemoryUsageAllSensorHandler)

	log.Println("server listening on :8080")
	go http.ListenAndServe("0.0.0.0:8080", nil)
	select {}
}
