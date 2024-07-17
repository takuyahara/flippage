package listener

import (
	"fmt"
	"time"

	"github.com/takuyahara/flippage/config"
	"github.com/takuyahara/flippage/utils"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

const TICK_INTERVAL_MSEC = 100

var chStop chan bool
var chReset chan bool

func Stop() {
	chStop <- false
	removeEventHooks()
}

func removeEventHooks() {
	hook.End()
}

func addEventHooks() {
	resetCounter := func(e hook.Event) {
		chReset <- true
	}
	hook.Register(hook.KeyDown, []string{`left`}, resetCounter)
	hook.Register(hook.KeyDown, []string{`right`}, resetCounter)
	hook.Register(hook.KeyDown, []string{`up`}, resetCounter)
	hook.Register(hook.KeyHold, []string{`down`}, resetCounter) // Workaround: Oddly, key down event for `down` will never be invoked
	hook.Register(hook.MouseWheel, []string{}, resetCounter)
	// hook.Register(hook.MouseDown, []string{}, resetCounter) // Seems not working
	hook.Register(hook.MouseUp, []string{}, resetCounter)
	hook.Register(hook.MouseDrag, []string{}, resetCounter)
	s := hook.Start()
	<-hook.Process(s)
}

func ListenEvents() {
	go addEventHooks()
}

func Flip(msg1 string, direction, interval int, key string) {
	cnt := 0.0
	incr := float64(TICK_INTERVAL_MSEC) / 1000
	duration := time.Millisecond * TICK_INTERVAL_MSEC
	// Main func
	mainFunc := func() {
		remaining := interval - int(cnt)
		if remaining <= 0 {
			removeEventHooks()
			robotgo.KeyTap(key)
			go addEventHooks()
			remaining = interval
			cnt = 0
		}
		progress := cnt / float64(interval)
		var msg2 string
		if remaining > 1 {
			msg2 = fmt.Sprintf(`Will flip page in %d seconds...`, remaining)
		} else {
			msg2 = `Will flip page in 1 second...`
		}
		message := msg1 + msg2
		if direction == config.DIRECTION_DOWN {
			utils.ShowProgressBarDown(progress, message)
		} else {
			if direction == config.DIRECTION_LEFT {
				utils.ShowProgressBarLeft(progress, message)
			} else {
				utils.ShowProgressBarRight(progress, message)
			}
		}
		cnt += incr
	}
	mainFunc() // Run once immediately
	// Run ticker
	chStop = make(chan bool, 1)
	chReset = make(chan bool, 1)
	ticker := time.NewTicker(duration)
tick:
	for {
		select {
		case <-ticker.C:
			mainFunc()
		case <-chReset:
			cnt = 0
			mainFunc()
			ticker.Reset(duration)
		case <-chStop:
			ticker.Stop()
			break tick
		}
	}
}

func Scroll(msg1 string, interval int) {
	cnt := 0.0
	incr := float64(TICK_INTERVAL_MSEC) / 1000
	duration := time.Millisecond * TICK_INTERVAL_MSEC
	// Main func
	mainFunc := func() {
		remaining := interval - int(cnt)
		if remaining <= 0 {
			removeEventHooks()
			robotgo.ScrollDir(1, `down`)
			go addEventHooks()
			remaining = interval
			cnt = 0
		}
		progress := cnt / float64(interval)
		var msg2 string
		if remaining > 1 {
			msg2 = fmt.Sprintf(`Will scroll page in %d seconds...`, remaining)
		} else {
			msg2 = `Will scroll page in 1 second...`
		}
		message := msg1 + msg2
		utils.ShowProgressBarScroll(progress, message)
		cnt += incr
	}
	mainFunc() // Run once immediately
	// Run ticker
	chStop = make(chan bool)
	chReset = make(chan bool)
	ticker := time.NewTicker(duration)
tick:
	for {
		select {
		case <-ticker.C:
			mainFunc()
		case <-chReset:
			cnt = 0
			mainFunc()
			ticker.Reset(duration)
		case <-chStop:
			ticker.Stop()
			break tick
		}
	}
}
