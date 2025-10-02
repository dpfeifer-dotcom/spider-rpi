package hardwares

import (
	"log"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type LedControl struct {
	Option     ws2811.Option
	WS2811     *ws2811.WS2811
	Running    bool
	Colors     Colors
	LedQueue   []int
	Leds       []uint32
	LedService LedService
}

type Colors struct {
	ORANGE1 uint32
	ORANGE2 uint32
	ORANGE3 uint32
	RED1    uint32
	RED2    uint32
	BLUE1   uint32
	BLUE2   uint32
	NONE    uint32
}

func NewLedControll() *LedControl {
	var ledError error

	ledControll := LedControl{Option: ws2811.DefaultOptions,
		Running: false,
		Colors: Colors{
			ORANGE1: RGBToColor(255, 160, 0),
			ORANGE2: RGBToColor(127, 80, 0),
			ORANGE3: RGBToColor(63, 40, 0),
			RED1:    RGBToColor(255, 0, 0),
			RED2:    RGBToColor(127, 0, 0),
			BLUE1:   RGBToColor(0, 0, 255),
			BLUE2:   RGBToColor(0, 0, 190),
			NONE:    RGBToColor(0, 0, 0)},
		LedQueue: []int{0, 1, 2, 5, 4, 3}}

	ledControll.setOptions()
	ledControll.WS2811, ledError = ws2811.MakeWS2811(&ledControll.Option)
	if ledError != nil {
		log.Println(ledError)
	}
	// Inicializáció
	if ledError = ledControll.WS2811.Init(); ledError != nil {
		log.Println(ledError)
	}
	ledControll.Leds = ledControll.WS2811.Leds(0)

	return &ledControll
}

func (ledControll *LedControl) SetAction(ledService LedService) {
	ledControll.LedService.Stop()

	for !ledControll.LedService.IsStopped() {
		log.Println("sleep")
		time.Sleep(100 * time.Millisecond)
	}

	ledControll.LedService = ledService
	ledControll.LedService.Start()
}

func (ledControll *LedControl) setOptions() {
	ledControll.Option.Channels[0].LedCount = 6
	ledControll.Option.Channels[0].GpioPin = 12
	ledControll.Option.Channels[0].Brightness = 255
	ledControll.Option.Channels[0].StripeType = ws2811.WS2811StripGRB
	ledControll.Option.Channels[0].Invert = false
}

func (ledControl *LedControl) TurnOff() {
	for i := range ledControl.LedQueue {
		ledControl.Leds[ledControl.LedQueue[i]] = ledControl.Colors.NONE
	}
	if err := ledControl.WS2811.Render(); err != nil {
		log.Fatalf("Error rendering LEDs: %v", err)
	}
}

type LedService interface {
	Start()
	Stop()
	IsStopped() bool
}

func RGBToColor(r, g, b uint8) uint32 {
	return uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}
