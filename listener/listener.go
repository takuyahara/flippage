package listener

import (
	"flippage/config"
	"flippage/utils"
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"

	"github.com/micmonay/keybd_event"
)

var cnt int
var chStop chan bool
var skipEvent bool

func Stop() {
	chStop <- false
	robotgo.EventEnd()
}

func ListenEvents() {
	skipEvent = false
	resetCounter := func(e hook.Event) {
		if !skipEvent {
			cnt = 0
		}
		skipEvent = false
	}
	robotgo.EventHook(hook.KeyDown, []string{`left`}, resetCounter)
	robotgo.EventHook(hook.KeyDown, []string{`right`}, resetCounter)
	robotgo.EventHook(hook.KeyDown, []string{`up`}, resetCounter)
	robotgo.EventHook(hook.KeyHold, []string{`down`}, resetCounter) // Workaround: Oddly, key down event for `down` will never be invoked
	robotgo.EventHook(hook.MouseWheel, []string{}, resetCounter)
	robotgo.EventHook(hook.MouseDown, []string{}, resetCounter)
	robotgo.EventHook(hook.MouseUp, []string{}, resetCounter)
	robotgo.EventHook(hook.MouseDrag, []string{}, resetCounter)
	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func Flip(msg1 string, direction, interval, vk int) {
	cnt = 0
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(vk)
	// Main func
	mainFunc := func() {
		remaining := interval - cnt
		if remaining <= 0 {
			skipEvent = true
			if err := kb.Launching(); err != nil {
				panic(err)
			}
			remaining = interval
			cnt = 0
		}
		progress := 1.0 - float64(remaining-1)/float64(interval)
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
		cnt++
	}
	mainFunc() // Run once immediately
	// Run ticker
	chStop = make(chan bool)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			mainFunc()
		case <-chStop:
			ticker.Stop()
			break
		}
	}
}

func Scroll(msg1 string, interval int) {
	cnt = 0
	// Main func
	mainFunc := func() {
		remaining := interval - cnt
		if remaining <= 0 {
			skipEvent = true
			robotgo.ScrollMouse(1, `down`)
			remaining = interval
			cnt = 0
		}
		progress := 1.0 - float64(remaining-1)/float64(interval)
		var msg2 string
		if remaining > 1 {
			msg2 = fmt.Sprintf(`Will scroll page in %d seconds...`, remaining)
		} else {
			msg2 = `Will scroll page in 1 second...`
		}
		message := msg1 + msg2
		utils.ShowProgressBarScroll(progress, message)
		cnt++
	}
	mainFunc() // Run once immediately
	// Run ticker
	chStop = make(chan bool)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			mainFunc()
		case <-chStop:
			ticker.Stop()
			break
		}
	}
}
