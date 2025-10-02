package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type MemoryListCache struct {
	Name         *string
	settings     settingsList
	cacheStorage []interface{}
	mutex        *sync.Mutex
	hasChange    *bool
	colorPalette Colors
}

type settingsList struct {
	DefaultExpirationSec        int
	ExpirationCheckDelaySeconds float64
	Limit                       int
	verbose                     bool
	persistent                  bool
}

func NewMemoryListCache(name string) *MemoryListCache {
	memoryCache := &MemoryListCache{
		Name: &name,
		settings: settingsList{
			ExpirationCheckDelaySeconds: 1,
			DefaultExpirationSec:        -1,
			Limit:                       0,
			verbose:                     false,
			persistent:                  false,
		},
		mutex:        &sync.Mutex{},
		cacheStorage: make([]any, 0),
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

func (memoryCache *MemoryListCache) SetLimits(limit int) *MemoryListCache {
	if limit > 0 {
		memoryCache.settings.Limit = limit
	}
	memoryCache.print("set limit", memoryCache.colorPalette.DEFAULT)

	return memoryCache
}

func (memoryCache *MemoryListCache) SetDefaultData(data any) *MemoryListCache {
	memoryCache.cacheStorage = make([]any, memoryCache.settings.Limit)
	for i := range memoryCache.cacheStorage {
		memoryCache.cacheStorage[i] = data
	}
	return memoryCache
}

func (memoryCache *MemoryListCache) StartWorkers() *MemoryListCache {
	if memoryCache.settings.persistent {
		go memoryCache.startPeristierCheckWorker()
	}

	return memoryCache
}
func (memoryCache *MemoryListCache) Verbose() *MemoryListCache {
	memoryCache.settings.verbose = true
	memoryCache.print("verbose", memoryCache.colorPalette.DEFAULT)
	return memoryCache
}

func (memoryCache *MemoryListCache) Persistent() *MemoryListCache {
	memoryCache.settings.persistent = true
	start := time.Now()
	memoryCache.print("start loading persited items", memoryCache.colorPalette.DEFAULT)
	readErr := memoryCache.readFromFile()
	if readErr != nil {
		memoryCache.printError(readErr.Error())
	}
	memoryCache.print("finished loading persited items: "+time.Since(start).String(), memoryCache.colorPalette.DEFAULT)
	memoryCache.print("presisted items: "+strconv.Itoa(len(memoryCache.cacheStorage)), memoryCache.colorPalette.DEFAULT)
	return memoryCache
}

func (memoryCache *MemoryListCache) Set(data interface{}) *MemoryListCache {
	memoryCache.mutex.Lock()
	if memoryCache.settings.Limit > 0 && len(memoryCache.cacheStorage) == memoryCache.settings.Limit {
		memoryCache.cacheStorage = append(memoryCache.cacheStorage[:0], memoryCache.cacheStorage[1:]...)
	}
	memoryCache.cacheStorage = append(memoryCache.cacheStorage, data)
	memoryCache.mutex.Unlock()
	memoryCache.print("set item", memoryCache.colorPalette.GREEN)
	return memoryCache
}

func (memoryCache *MemoryListCache) GetLast() any {
	memoryCache.mutex.Lock()
	var cachedItem any
	if len(memoryCache.cacheStorage) > 0 {
		cachedItem = memoryCache.cacheStorage[len(memoryCache.cacheStorage)-1]
	}
	memoryCache.mutex.Unlock()
	memoryCache.print("get item", memoryCache.colorPalette.GREEN)
	return cachedItem
}

func (memoryCache *MemoryListCache) GetAll() []any {
	memoryCache.mutex.Lock()
	cachedItems := make([]any, 0)
	if len(memoryCache.cacheStorage) > 0 {
		cachedItems = memoryCache.cacheStorage[:len(memoryCache.cacheStorage)-1]
	}
	memoryCache.mutex.Unlock()
	memoryCache.print("get all items", memoryCache.colorPalette.GREEN)
	return cachedItems
}

func (memoryCache *MemoryListCache) print(text string, color string) {
	if memoryCache.settings.verbose {
		log.Println(memoryCache.colorPalette.CYAN +
			"[" + *memoryCache.Name + "-cache]" +
			color +
			" " +
			text +
			memoryCache.colorPalette.DEFAULT)
	}
}

func (memoryCache *MemoryListCache) printError(text string) {
	log.Println(memoryCache.colorPalette.CYAN +
		"[" + *memoryCache.Name + "-cache]" +
		memoryCache.colorPalette.RED +
		" " +
		text +
		memoryCache.colorPalette.DEFAULT)
}

func (memoryCache *MemoryListCache) startPeristierCheckWorker() {
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
func (memoryCache *MemoryListCache) writeToFile() error {
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

func (memoryCache *MemoryListCache) readFromFile() error {
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
