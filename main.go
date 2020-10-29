package main

import (
	"bufio"
	"flippage/appinfo"
	"flippage/keyboard"
	"flippage/utils"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/micmonay/keybd_event"
)

const APP_NAME = `Flippage`

func main() {
	var interval int
	var vk int
	utils.Clear()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(`Type interval in second which is positive non-zero int: `)
		scanner.Scan()
		scanned := scanner.Text()
		if regexp.MustCompile(`^[1-9][0-9]*$`).MatchString(scanned) {
			i, err := strconv.Atoi(scanned)
			if err != nil {
				panic(err)
			}
			interval = i
			break
		}
		utils.ClearPreviousLine()
	}
	for {
		utils.ClearPreviousLine()
		fmt.Print(`Specify a direction to flip (left, right, up, down): `)
		scanner.Scan()
		scanned := scanner.Text()
		if regexp.MustCompile(`^(?:left|right|up|down)$`).MatchString(scanned) {
			switch scanned {
			case `left`:
				vk = keybd_event.VK_LEFT
			case `right`:
				vk = keybd_event.VK_RIGHT
			case `up`:
				vk = keybd_event.VK_UP
			case `down`:
				vk = keybd_event.VK_DOWN
			}
			break
		}
		utils.ClearLine()
	}
	chAppInfoForeground := make(chan appinfo.AppInfo, 2)
	go keyboard.ListenEvents()
	// Wait until foreground app changes
	utils.ClearPreviousLine()
	fmt.Print(`Waiting to foreground app changes...`)
	appinfo.GetChangedForegroundInfo(chAppInfoForeground)
	appInfoTarget := <-chAppInfoForeground
	// Flip page automatically
	utils.ClearLine()
	fmt.Printf("%s has activated for %s.\n", APP_NAME, appInfoTarget.Name)
	go keyboard.Send(interval, vk)
	appinfo.GetChangedForegroundInfo(chAppInfoForeground)
	<-chAppInfoForeground
	// Close app
	utils.ClearLine()
	fmt.Printf("%s has exited as foreground app has changed.\n", APP_NAME)
}
