package main

import (
	"fmt"
	"log"
	"net/http"
	"spider-sensor/globals"
	"spider-sensor/handlers"
	"spider-sensor/hardwares"
	cache "spider-sensor/memorycache"
	"spider-sensor/services"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

const (
	MPU6050_ADDR         = 0x69
	MPU6050_PWR_MGMT_1   = 0x6B
	MPU6050_ACCEL_XOUT_H = 0x3B
)

func readWord16(dev i2c.Dev, reg byte) (int16, error) {
	buf := []byte{0, 0}
	if err := dev.Tx([]byte{reg}, buf); err != nil {
		return 0, err
	}
	val := int16(buf[0])<<8 | int16(buf[1])
	return val, nil
}

func main() {
	// Inicializáljuk a periph.io-t
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Megnyitjuk az I2C buszt
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// Létrehozzuk az MPU6050 eszközt
	dev := i2c.Dev{Addr: MPU6050_ADDR, Bus: bus}

	// Kioldjuk az alvó módot
	if err := dev.Tx([]byte{MPU6050_PWR_MGMT_1, 0}, nil); err != nil {
		log.Fatal("Failed to wake up MPU6050:", err)
	}
	go func() {
		for {
			ax, err := readWord16(dev, MPU6050_ACCEL_XOUT_H)
			if err != nil {
				log.Println("Error reading acceleration:", err)
				continue
			}
			ay, _ := readWord16(dev, MPU6050_ACCEL_XOUT_H+2)
			az, _ := readWord16(dev, MPU6050_ACCEL_XOUT_H+4)

			fmt.Printf("Accel X: %d, Y: %d, Z: %d\n", ax, ay, az)
			time.Sleep(500 * time.Millisecond)
		}
	}()

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
