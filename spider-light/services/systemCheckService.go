package services

import (
	cache "spider-light/memorycache"
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
