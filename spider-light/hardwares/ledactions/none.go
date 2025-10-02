package ledactions

import "spider-light/hardwares"

type None struct {
	ledControl *hardwares.LedControl
	running    bool
	stopped    bool
}

func NoneLight(ledcontroll *hardwares.LedControl) *None {
	noneLight := None{
		ledControl: ledcontroll,
		running:    false,
		stopped:    true}
	return &noneLight
}

func (none *None) IsStopped() bool {
	return none.stopped
}

func (none *None) Start() {

}

func (none *None) Stop() {

}
