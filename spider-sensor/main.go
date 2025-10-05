package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"spider-sensor/globals"
	"spider-sensor/handlers"
	"spider-sensor/hardwares"
	cache "spider-sensor/memorycache"
	"spider-sensor/services"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	mpu6050 := i2c.NewMPU6050Driver(firmataAdaptor)

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			mpu6050.GetData()

			fmt.Println("Accelerometer", mpu6050.Accelerometer)
			fmt.Println("Gyroscope", mpu6050.Gyroscope)
			fmt.Println("Temperature", mpu6050.Temperature)
		})
	}

	robot := gobot.NewRobot("mpu6050Bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{mpu6050},
		work,
	)

	go robot.Start()

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
