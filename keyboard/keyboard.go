package keyboard

import (
	hook "flippage/gohook"
	"flippage/utils"
	"fmt"
	"time"

	"github.com/micmonay/keybd_event"
)

var cnt int

func ListenEvents() {
	hook.Register(hook.KeyDown, []string{`left`}, func(e hook.Event) {
		cnt = 0
	})
	hook.Register(hook.KeyDown, []string{`right`}, func(e hook.Event) {
		cnt = 0
	})
	hook.Register(hook.KeyDown, []string{`up`}, func(e hook.Event) {
		cnt = 0
	})
	// Workaround: Oddly, key down event for `down` will never be invoked
	hook.Register(hook.KeyHold, []string{`down`}, func(e hook.Event) {
		cnt = 0
	})
	s := hook.Start()
	<-hook.Process(s)
}

func Send(interval int, vk int) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(vk)
	cnt = 0
	for {
		if cnt >= interval {
			if err := kb.Launching(); err != nil {
				panic(err)
			}
			cnt = 0
		}
		utils.ClearLine()
		remaining := interval - cnt
		if remaining > 1 {
			fmt.Printf("Will flip page in %d seconds...", remaining)
		} else {
			fmt.Print(`Will flip page in 1 second...`)
		}
		cnt++
		time.Sleep(time.Second)
	}
}
