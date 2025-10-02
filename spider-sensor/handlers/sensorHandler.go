package handlers

import (
	"encoding/json"
	"net/http"
	"spider-sensor/globals"
)

func CPUSUsageSensorHandler(w http.ResponseWriter, r *http.Request) {

	sensorValue := globals.CPUUsageCache.GetLast().(float64)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]float64{"value": sensorValue})
}

func CPUSUsageAllSensorHandler(w http.ResponseWriter, r *http.Request) {

	sensorValue := globals.CPUUsageCache.GetAll()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string][]any{"value": sensorValue})
}

func CPUSTempSensorHandler(w http.ResponseWriter, r *http.Request) {

	sensorValue := globals.CPUTempCache.GetLast().(float64)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]float64{"value": sensorValue})
}

func CPUSTempAllSensorHandler(w http.ResponseWriter, r *http.Request) {

	sensorValue := globals.CPUTempCache.GetAll()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string][]any{"value": sensorValue})
}

func MemoryUsageSensorHandler(w http.ResponseWriter, r *http.Request) {

	sensorValue := globals.MemoryUsageCache.GetLast().(float64)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]float64{"value": sensorValue})
}

func MemoryUsageAllSensorHandler(w http.ResponseWriter, r *http.Request) {

	sensorValue := globals.MemoryUsageCache.GetAll()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string][]any{"value": sensorValue})
}
