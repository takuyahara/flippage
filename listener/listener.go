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
var isRunning bool
var direction int
var mode int

func Stop() {
	isRunning = false
	robotgo.EventEnd()
	// if mode == config.MODE_SCROLL || direction == config.DIRECTION_DOWN {
	// 	utils.ClearVerticalProgressBar()
	// } else {
	// 	utils.ClearHorizontalProgressBar()
	// }
}

func ListenEvents() {
	robotgo.EventHook(hook.KeyDown, []string{`left`}, func(e hook.Event) {
		cnt = 0
	})
	robotgo.EventHook(hook.KeyDown, []string{`right`}, func(e hook.Event) {
		cnt = 0
	})
	robotgo.EventHook(hook.KeyDown, []string{`up`}, func(e hook.Event) {
		cnt = 0
	})
	// Workaround: Oddly, key down event for `down` will never be invoked
	robotgo.EventHook(hook.KeyHold, []string{`down`}, func(e hook.Event) {
		cnt = 0
	})
	robotgo.EventHook(hook.MouseWheel, []string{}, func(e hook.Event) {
		cnt = 0
	})
	robotgo.EventHook(hook.MouseDown, []string{}, func(e hook.Event) {
		cnt = 0
	})
	robotgo.EventHook(hook.MouseUp, []string{}, func(e hook.Event) {
		cnt = 0
	})
	robotgo.EventHook(hook.MouseDrag, []string{}, func(e hook.Event) {
		cnt = 0
	})
	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func Flip(msg1 string, d, interval, vk int) {
	mode = config.MODE_FLIP
	direction = d
	isRunning = true
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(vk)
	cnt = 0
	for isRunning {
		remaining := interval - cnt
		if remaining <= 0 {
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
		time.Sleep(time.Second)
	}
}

func Scroll(msg1 string, interval int) {
	mode = config.MODE_SCROLL
	isRunning = true
	cnt = 0
	for isRunning {
		remaining := interval - cnt
		if remaining <= 0 {
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
		time.Sleep(time.Second)
	}
}
