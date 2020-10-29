package listener

import (
	"flippage/utils"
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"

	"github.com/micmonay/keybd_event"
)

var cnt int
var isRunning bool

func Stop() {
	isRunning = false
	robotgo.EventEnd()
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

func Flip(interval int, vk int) {
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
		utils.ClearLine()
		if remaining > 1 {
			fmt.Printf("Will flip page in %d seconds...", remaining)
		} else {
			fmt.Printf("Will flip page in 1 second...")
		}
		cnt++
		time.Sleep(time.Second)
	}
}

func Scroll(interval int) {
	isRunning = true
	cnt = 0
	for isRunning {
		remaining := interval - cnt
		if remaining <= 0 {
			robotgo.ScrollMouse(1, `down`)
			remaining = interval
			cnt = 0
		}
		utils.ClearLine()
		if remaining > 1 {
			fmt.Printf("Will scroll page in %d seconds...", remaining)
		} else {
			fmt.Printf("Will scroll page in 1 second...")
		}
		cnt++
		time.Sleep(time.Second)
	}
}
