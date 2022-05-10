package platform

import (
	"time"

	"tractor.dev/apptron/bridge/platform/win32"
)

var mainfunc = make(chan func())
var quit = make(chan bool)

func Main() {
	win32.OS_Init()

	runReady.Wait()

	// @Robustness: should we add some sort of framerate timing to this?
loop:
	for {
		win32.PollEvents()

		select {
		case fn := <-mainfunc:
			fn()
		case <-quit:
			break loop
		default: // NOTE(nick): keep running at max speed!
		}

		// @Robustness: how accurate is time.Sleep? should we use win32.SleepMS instead?
		time.Sleep(1 * time.Millisecond)
	}
}

func Terminate() {
	quit <- true
}