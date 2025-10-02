package hardwares

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetCpuUsage() float64 {
	percent, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return 0
	}
	return percent[0]
}

func GetCpuTemperature() float64 {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return 0
	}
	return temps[0].Temperature
}

func GetMemoryUsage() float64 {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return vmStat.UsedPercent
}
