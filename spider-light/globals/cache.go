package globals

import (
	"spider-light/hardwares"

	"github.com/orsinium-labs/gamepad"
	"periph.io/x/conn/v3/gpio"
)

var Controller *gamepad.GamePad
var LedControl *hardwares.LedControl
var Bulb gpio.PinIO
