package hardwares

import (
	"log"

	"github.com/orsinium-labs/gamepad"
)

func NewController(gamepadId int) (controller *gamepad.GamePad) {
	controller, controllerError := gamepad.NewGamepad(gamepadId)
	if controllerError != nil {
		log.Println(controllerError)
	}
	return controller
}

func ButtonPressed(state bool, previousState *bool, command func()) {
	if state {
		if !*previousState {
			*previousState = true
			command()

		}
	} else {
		if *previousState {
		}
		*previousState = false
	}

}
