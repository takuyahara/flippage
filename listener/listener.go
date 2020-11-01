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

var chStop chan bool
var chReset chan bool

func Stop() {
	chStop <- false
	removeEventHooks()
}

func removeEventHooks() {
	robotgo.EventEnd()
}

func addEventHooks() {
	resetCounter := func(e hook.Event) {
		chReset <- true
	}
	robotgo.EventHook(hook.KeyDown, []string{`left`}, resetCounter)
	robotgo.EventHook(hook.KeyDown, []string{`right`}, resetCounter)
	robotgo.EventHook(hook.KeyDown, []string{`up`}, resetCounter)
	robotgo.EventHook(hook.KeyHold, []string{`down`}, resetCounter) // Workaround: Oddly, key down event for `down` will never be invoked
	robotgo.EventHook(hook.MouseWheel, []string{}, resetCounter)
	// robotgo.EventHook(hook.MouseDown, []string{}, resetCounter) // Seems not working
	robotgo.EventHook(hook.MouseUp, []string{}, resetCounter)
	robotgo.EventHook(hook.MouseDrag, []string{}, resetCounter)
	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func ListenEvents() {
	go addEventHooks()
}

func Flip(msg1 string, direction, interval, vk int) {
	cnt := 0
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(vk)
	// Main func
	mainFunc := func() {
		remaining := interval - cnt
		if remaining <= 0 {
			removeEventHooks()
			if err := kb.Launching(); err != nil {
				panic(err)
			}
			go addEventHooks()
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
	chStop = make(chan bool, 1)
	chReset = make(chan bool, 1)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			mainFunc()
		case <-chReset:
			cnt = 0
			mainFunc()
			ticker.Reset(time.Second)
		case <-chStop:
			ticker.Stop()
			break
		}
	}
}

func Scroll(msg1 string, interval int) {
	cnt := 0
	// Main func
	mainFunc := func() {
		remaining := interval - cnt
		if remaining <= 0 {
			removeEventHooks()
			robotgo.ScrollMouse(1, `down`)
			go addEventHooks()
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
	chReset = make(chan bool)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			mainFunc()
		case <-chReset:
			cnt = 0
			mainFunc()
			ticker.Reset(time.Second)
		case <-chStop:
			ticker.Stop()
			break
		}
	}
}
