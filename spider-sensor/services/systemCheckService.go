package services

import (
	cache "spider-sensor/memorycache"
	"time"
)

func SystemListCheckService(cache *cache.MemoryListCache, getter func() float64, sleep int) {
	go func() {
		for {
			cache.Set(getter())
			time.Sleep(time.Duration(sleep) * time.Second)
		}
	}()
}
