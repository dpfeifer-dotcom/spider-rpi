package globals

import (
	"spider-backend/hardwares"
	cache "spider-backend/memorycache"

	"github.com/orsinium-labs/gamepad"
	"periph.io/x/conn/v3/gpio"
)

var SystemCache *cache.MemoryCache
var CPUTempCache *cache.MemoryListCache
var CPUUsageCache *cache.MemoryListCache
var MemoryUsageCache *cache.MemoryListCache
var Controller *gamepad.GamePad
var LedControl *hardwares.LedControl
var Bulb gpio.PinIO
