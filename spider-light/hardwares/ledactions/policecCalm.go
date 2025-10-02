package ledactions

import (
	"log"
	"spider-light/hardwares"
	"time"
)

type PoliceCalm struct {
	ledControl *hardwares.LedControl
	running    bool
	stopped    bool
}

func PoliceCalmLight(ledControl *hardwares.LedControl) *PoliceCalm {
	policeCalmLight := PoliceCalm{
		ledControl: ledControl,
		running:    false,
		stopped:    true,
	}
	return &policeCalmLight
}

func (policeCalm *PoliceCalm) IsStopped() bool {
	return policeCalm.stopped
}

func (policeCalm *PoliceCalm) Start() {
	policeCalm.running = true
	policeCalm.stopped = false
	go func() {
		leds := policeCalm.ledControl.Leds
		for policeCalm.running {
			for _, ledId := range policeCalm.ledControl.LedQueue {
				for i := range policeCalm.ledControl.LedQueue {
					leds[policeCalm.ledControl.LedQueue[i]] = policeCalm.ledControl.Colors.BLUE2
				}
				for i := range policeCalm.ledControl.LedQueue {
					if ledId == policeCalm.ledControl.LedQueue[i] {
						leds[policeCalm.ledControl.LedQueue[i]] = policeCalm.ledControl.Colors.RED1
						if i == 0 {
							leds[policeCalm.ledControl.LedQueue[5]] = policeCalm.ledControl.Colors.RED2
						} else {
							leds[policeCalm.ledControl.LedQueue[i-1]] = policeCalm.ledControl.Colors.RED2
						}
					}
				}
				if err := policeCalm.ledControl.WS2811.Render(); err != nil {
					log.Fatalf("Error rendering LEDs: %v", err)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
		policeCalm.ledControl.TurnOff()
		policeCalm.stopped = true
	}()
}

func (policeCalm *PoliceCalm) Stop() {
	policeCalm.running = false
}
