package hardwares

import (
	"fmt"

	"periph.io/x/conn/v3/gpio"
)

func SwitchBulbLight(pin gpio.PinIO) {
	fmt.Printf("pin.Read(): %v\n", pin.Read())
	if pin.Read() == gpio.Low {
		pin.Out(gpio.High)
	} else {
		pin.Out(gpio.Low)
	}
}
