package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

type MemoryCache struct {
	Name         *string
	settings     settings
	cacheStorage map[string]*CachedItem
	mutex        *sync.Mutex
	hasChange    *bool
	colorPalette Colors
}

type settings struct {
	DefaultExpirationSec        int
	ExpirationCheckDelaySeconds float64
	verbose                     bool
	persistent                  bool
}

type CachedItem struct {
	Value       interface{}
	CreatedAt   time.Time
	ExipritySec int
}

type CacheResult struct {
	Key        string
	Value      interface{}
	Exists     bool
	Expiration int
}

type Colors struct {
	DEFAULT string
	BLACK   string
	RED     string
	GREEN   string
	Yellow  string
	BLUE    string
	MAGENTA string
	CYAN    string
	WHITE   string
}

func NewMemoryCache(name string) *MemoryCache {
	memoryCache := &MemoryCache{
		Name: &name,
		settings: settings{
			ExpirationCheckDelaySeconds: 1,
			DefaultExpirationSec:        -1,
			verbose:                     false,
			persistent:                  false,
		},
		mutex:        &sync.Mutex{},
		cacheStorage: make(map[string]*CachedItem),
		hasChange:    &[]bool{false}[0],
		colorPalette: Colors{
			DEFAULT: "\033[0m",
			BLACK:   "\033[30m",
			RED:     "\033[31m",
			GREEN:   "\033[32m",
			Yellow:  "\033[33m",
			BLUE:    "\033[34m",
			MAGENTA: "\033[35m",
			CYAN:    "\033[36m",
			WHITE:   "\033[37m",
		},
	}
	return memoryCache
}

func (memoryCache *MemoryCache) StartWorkers() *MemoryCache {
	go memoryCache.startExpirationCheckWorker()
	if memoryCache.settings.persistent {
		go memoryCache.startPeristierCheckWorker()
	}

	return memoryCache
}
func (memoryCache *MemoryCache) Verbose() *MemoryCache {
	memoryCache.settings.verbose = true
	memoryCache.print("verbose", memoryCache.colorPalette.DEFAULT)
	return memoryCache
}

func (memoryCache *MemoryCache) Persistent() *MemoryCache {
	memoryCache.settings.persistent = true
	start := time.Now()
	memoryCache.print("start loading persited items", memoryCache.colorPalette.DEFAULT)
	readErr := memoryCache.readFromFile()
	if readErr != nil {
		memoryCache.printError(readErr.Error())
	}
	memoryCache.print("finished loading persited items: "+time.Since(start).String(), memoryCache.colorPalette.DEFAULT)
	memoryCache.print("presisted items: "+strconv.Itoa(len(memoryCache.cacheStorage)), memoryCache.colorPalette.DEFAULT)
	deletedItems := memoryCache.checkExipredItems()
	memoryCache.print("deleted items: "+strconv.Itoa(deletedItems), memoryCache.colorPalette.DEFAULT)
	return memoryCache
}

func (memoryCache *MemoryCache) Set(key string, value interface{}, expiration ...int) {
	memoryCache.mutex.Lock()
	setExpiration := memoryCache.settings.DefaultExpirationSec
	if expiration != nil {
		setExpiration = expiration[0]
	}
	memoryCache.cacheStorage[key] = &CachedItem{
		Value:       value,
		CreatedAt:   time.Now(),
		ExipritySec: setExpiration}
	*memoryCache.hasChange = true
	memoryCache.mutex.Unlock()
	memoryCache.print("set item: "+key, memoryCache.colorPalette.GREEN)
}

func (memoryCache *MemoryCache) Get(key string) (result CacheResult) {
	memoryCache.mutex.Lock()
	value, exists := memoryCache.cacheStorage[key]
	memoryCache.mutex.Unlock()
	if exists {
		return CacheResult{
			Key:        key,
			Value:      value.Value,
			Exists:     exists,
			Expiration: value.ExipritySec,
		}
	}
	memoryCache.print("get item: "+key, memoryCache.colorPalette.DEFAULT)
	return CacheResult{}
}

func (memoryCache *MemoryCache) Mod(key string, value interface{}) {
	memoryCache.mutex.Lock()
	_, exists := memoryCache.cacheStorage[key]
	if exists {
		memoryCache.cacheStorage[key].Value = value
	}
	*memoryCache.hasChange = true
	memoryCache.mutex.Unlock()
}

func (memoryCache *MemoryCache) Del(key string) {
	memoryCache.mutex.Lock()
	delete(memoryCache.cacheStorage, key)
	*memoryCache.hasChange = true
	memoryCache.mutex.Unlock()
	memoryCache.print("delete item: "+key, memoryCache.colorPalette.DEFAULT)
}

func (memoryCache *MemoryCache) Keys() (keys []string) {
	memoryCache.mutex.Lock()
	for key := range memoryCache.cacheStorage {
		keys = append(keys, key)
	}
	memoryCache.mutex.Unlock()
	memoryCache.print("get all keys", memoryCache.colorPalette.DEFAULT)
	return
}

func (memoryCache *MemoryCache) print(text string, color string) {
	if memoryCache.settings.verbose {
		log.Println(memoryCache.colorPalette.CYAN +
			"[" + *memoryCache.Name + "-cache]" +
			color +
			" " +
			text +
			memoryCache.colorPalette.DEFAULT)
	}
}

func (memoryCache *MemoryCache) printError(text string) {
	log.Println(memoryCache.colorPalette.CYAN +
		"[" + *memoryCache.Name + "-cache]" +
		memoryCache.colorPalette.RED +
		" " +
		text +
		memoryCache.colorPalette.DEFAULT)
}

func (memoryCache *MemoryCache) DefaultExpirationSec(defaultExpirationSec int) *MemoryCache {
	memoryCache.settings.DefaultExpirationSec = defaultExpirationSec
	memoryCache.print("set default item expiration: "+strconv.Itoa(defaultExpirationSec)+"s", memoryCache.colorPalette.DEFAULT)
	return memoryCache
}

func (memoryCache *MemoryCache) ExpirationCheckDelaySeconds(checkDelaySeconds float64) *MemoryCache {
	memoryCache.settings.ExpirationCheckDelaySeconds = checkDelaySeconds
	stringCheckDelaySeconds := fmt.Sprintf("%.0f", checkDelaySeconds)
	memoryCache.print("set item expiration delay: "+stringCheckDelaySeconds+"s", memoryCache.colorPalette.DEFAULT)

	return memoryCache
}

func (memoryCache *MemoryCache) startExpirationCheckWorker() {
	expirationModifier := memoryCache.settings.ExpirationCheckDelaySeconds
	memoryCache.print("start expiration check worker", memoryCache.colorPalette.DEFAULT)
	for {
		start := time.Now()
		memoryCache.mutex.Lock()
		for itemKey := range memoryCache.cacheStorage {
			storedItem := memoryCache.cacheStorage[itemKey]
			if storedItem.ExipritySec != -1 {
				if storedItem.ExipritySec == 0 {
					delete(memoryCache.cacheStorage, itemKey)
					*memoryCache.hasChange = true
					memoryCache.print("delete expired item: "+itemKey, memoryCache.colorPalette.DEFAULT)
				} else {
					if storedItem.ExipritySec > int(expirationModifier) {
						storedItem.ExipritySec -= int(expirationModifier)
					} else {
						storedItem.ExipritySec = 0
					}
				}
			}
		}
		memoryCache.mutex.Unlock()
		latency := time.Since(start).Seconds()
		if latency <= memoryCache.settings.ExpirationCheckDelaySeconds {
			expirationModifier = memoryCache.settings.ExpirationCheckDelaySeconds
		} else {
			expirationModifier = math.Ceil(latency)
		}
		sleepTime := expirationModifier - latency
		time.Sleep(time.Duration(sleepTime * float64(time.Second)))
	}
}
func (memoryCache *MemoryCache) startPeristierCheckWorker() {
	for {
		time.Sleep(5 * time.Second)
		memoryCache.mutex.Lock()
		if *memoryCache.hasChange {
			memoryCache.print("cache change", memoryCache.colorPalette.DEFAULT)
			start := time.Now()
			err := memoryCache.writeToFile()
			elapsed := time.Since(start)
			memoryCache.print("save change: "+elapsed.String(), memoryCache.colorPalette.DEFAULT)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			*memoryCache.hasChange = false
		}
		memoryCache.mutex.Unlock()

	}
}
func (memoryCache *MemoryCache) writeToFile() error {
	file, err := os.Create(*memoryCache.Name + ".bin")
	if err != nil {
		return err
	}
	defer file.Close()
	gob.Register(map[string]*CachedItem{})
	gob.Register(CachedItem{})
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(memoryCache.cacheStorage)
	if err != nil {
		return err
	}
	return nil
}

func (memoryCache *MemoryCache) readFromFile() error {
	file, err := os.OpenFile(*memoryCache.Name+".bin", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	memoryCache.mutex.Lock()
	err = decoder.Decode(&memoryCache.cacheStorage)
	if err != nil {
		memoryCache.mutex.Unlock()
		if err == io.EOF {
			return nil
		}
		return err
	}
	memoryCache.mutex.Unlock()

	return nil
}

func (memoryCache *MemoryCache) checkExipredItems() (deletedItem int) {
	for key, item := range memoryCache.cacheStorage {
		if item.ExipritySec != -1 {
			item.ExipritySec -= int(time.Since(item.CreatedAt).Seconds())
			if item.ExipritySec < 0 {
				delete(memoryCache.cacheStorage, key)
				*memoryCache.hasChange = true
				deletedItem += 1
			}
		}
	}
	return
}
