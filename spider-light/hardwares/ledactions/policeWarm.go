package ledactions

import (
	"log"
	"spider-light/hardwares"
	"time"
)

type PoliceWarm struct {
	ledControl *hardwares.LedControl
	running    bool
	stopped    bool
}

func PoliceWarmlight(ledControl *hardwares.LedControl) *PoliceWarm {
	policeWarmlight := PoliceWarm{
		ledControl: ledControl,
		running:    false,
		stopped:    true,
	}
	return &policeWarmlight
}

func (policeWarm *PoliceWarm) IsStopped() bool {
	return policeWarm.stopped
}

func (policeWarm *PoliceWarm) Start() {
	toggle := true
	policeWarm.running = true
	policeWarm.stopped = false
	go func() {
		for policeWarm.running {
			if toggle {
				for i := 0; i < 3; i++ {
					policeWarm.ledControl.Leds[i] = policeWarm.ledControl.Colors.RED1
				}
				for i := 3; i < 6; i++ {
					policeWarm.ledControl.Leds[i] = policeWarm.ledControl.Colors.BLUE1
				}
			} else {
				for i := 0; i < 3; i++ {
					policeWarm.ledControl.Leds[i] = policeWarm.ledControl.Colors.BLUE1
				}
				for i := 3; i < 6; i++ {
					policeWarm.ledControl.Leds[i] = policeWarm.ledControl.Colors.RED1
				}
			}

			if err := policeWarm.ledControl.WS2811.Render(); err != nil {
				log.Fatalf("Error rendering LEDs: %v", err)
			}

			toggle = !toggle
			time.Sleep(350 * time.Millisecond)
		}
		policeWarm.ledControl.TurnOff()
		policeWarm.stopped = true
	}()
}

func (policeWarm *PoliceWarm) Stop() {
	policeWarm.running = false
}
