package hardwares

import "periph.io/x/conn/v3/gpio"

func SwitchBulbLight(pin gpio.PinIO) {
	if pin.Read() == gpio.Low {
		pin.Out(gpio.High)
	} else {
		pin.Out(gpio.Low)
	}
}
