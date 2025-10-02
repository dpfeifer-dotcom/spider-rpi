package services

import (
	cache "spider-backend/memorycache"
	"time"
)

func SystemCheckService(cache *cache.MemoryCache, name string, getter func() float64, sleep int) {
	go func() {
		for {
			cache.Set(name, getter())
			time.Sleep(time.Duration(sleep) * time.Second)
		}
	}()
}

func SystemListCheckService(cache *cache.MemoryListCache, getter func() float64, sleep int) {
	go func() {
		for {
			cache.Set(getter())
			time.Sleep(time.Duration(sleep) * time.Second)
		}
	}()
}
