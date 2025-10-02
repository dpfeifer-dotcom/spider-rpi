package ledactions

import (
	"log"
	"spider-light/hardwares"
	"time"
)

type Warning struct {
	ledControl *hardwares.LedControl
	running    bool
	stopped    bool
}

func WarningLight(ledcontroll *hardwares.LedControl) *Warning {
	warningLight := Warning{
		ledControl: ledcontroll,
		running:    false,
		stopped:    true}
	return &warningLight
}
func (warning *Warning) IsStopped() bool {
	return warning.stopped
}
func (warning *Warning) Start() {
	go func() {
		warning.running = true
		warning.stopped = false
		leds := warning.ledControl.Leds
		for warning.running {
			for _, ledId := range warning.ledControl.LedQueue {
				for i := range warning.ledControl.LedQueue {
					leds[warning.ledControl.LedQueue[i]] = warning.ledControl.Colors.ORANGE3
				}
				for i := range warning.ledControl.LedQueue {
					if ledId == warning.ledControl.LedQueue[i] {
						leds[warning.ledControl.LedQueue[i]] = warning.ledControl.Colors.ORANGE1
						if i == 0 {
							leds[warning.ledControl.LedQueue[5]] = warning.ledControl.Colors.ORANGE2
						} else {
							leds[warning.ledControl.LedQueue[i-1]] = warning.ledControl.Colors.ORANGE2
						}
					}
				}
				if err := warning.ledControl.WS2811.Render(); err != nil {
					warning.stopped = true
					log.Fatalf("Error rendering LEDs: %v", err)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
		warning.ledControl.TurnOff()
		warning.stopped = true
	}()
}

func (warning *Warning) Stop() {
	warning.running = false

}
