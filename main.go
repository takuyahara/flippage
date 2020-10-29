package main

import (
	"bufio"
	"flippage/appinfo"
	"flippage/listener"
	"flippage/utils"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/micmonay/keybd_event"
)

const APP_NAME = `Flippage`
const (
	RETRY_WITH_SAME_CONFIG = iota
	RETRY_WITH_DIFFERENT_CONFIG
	NO_RETRY
)

func getInterval() float64 {
	var interval float64
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(`Type interval in second which is positive float (%.1f): `)
		scanner.Scan()
		scanned := scanner.Text()
		if regexp.MustCompile(`^(?:0\.[1-9]|[1-9][0-9]*(?:\.[0-9])?)$`).MatchString(scanned) {
			i, err := strconv.ParseFloat(scanned, 64)
			if err != nil {
				panic(err)
			}
			interval = i
			break
		}
		utils.ClearPreviousLine()
	}
	return interval
}

func getVk() int {
	var vk int
	scanner := bufio.NewScanner(os.Stdin)
	for {
		utils.ClearPreviousLine()
		fmt.Print("Specify a direction to flip (\033[4ml\033[24meft/\033[4mr\033[24might/\033[4mu\033[24mp/\033[4md\033[24mown): ")
		scanner.Scan()
		scanned := scanner.Text()
		if regexp.MustCompile(`^(?:left|right|up|down|l|r|u|d)$`).MatchString(scanned) {
			switch scanned {
			case `l`, `left`:
				vk = keybd_event.VK_LEFT
			case `r`, `right`:
				vk = keybd_event.VK_RIGHT
			case `u`, `up`:
				vk = keybd_event.VK_UP
			case `d`, `down`:
				vk = keybd_event.VK_DOWN
			}
			break
		}
		utils.ClearLine()
	}
	return vk
}

func getRetry() uint {
	var retry uint
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(`Do you want to retry? (S/d/q)? 
 [S] Retry with same config
 [d] Retry with different config
 [q] Don't retry and quit app`)
		utils.GoToPreviousLine(3)
		fmt.Print(`Do you want to retry? (S/d/q)? `)
		scanner.Scan()
		scanned := scanner.Text()
		if regexp.MustCompile(`^(?:S|d|q)?$`).MatchString(scanned) {
			switch scanned {
			case ``:
			case `S`:
				retry = RETRY_WITH_SAME_CONFIG
			case `d`:
				retry = RETRY_WITH_DIFFERENT_CONFIG
			case `q`:
				retry = NO_RETRY
			}
			break
		}
		utils.GoToNextLine(3)
		utils.ClearPreviousLine(4)
	}
	utils.GoToNextLine(3)
	return retry
}

func run(interval float64, vk int) uint {
	chAppInfoForeground := make(chan appinfo.AppInfo, 2)
	go listener.ListenEvents()
	// Wait until foreground app changes
	utils.ClearPreviousLine()
	fmt.Print(`Waiting to foreground app changes...`)
	appinfo.GetChangedForegroundInfo(chAppInfoForeground)
	appInfoTarget := <-chAppInfoForeground
	// Flip page automatically
	utils.ClearLine()
	fmt.Printf("%s has activated for %s.\n", APP_NAME, appInfoTarget.Name)
	go listener.Flip(interval, vk)
	appinfo.GetChangedForegroundInfo(chAppInfoForeground)
	<-chAppInfoForeground
	// Close app
	utils.ClearLine()
	listener.Stop()
	fmt.Printf("%s has exited as foreground app has changed.\n", APP_NAME)
	return getRetry()
}

func main() {
	utils.Clear()
	interval := getInterval()
	vk := getVk()
	retry := run(interval, vk)
	for retry != NO_RETRY {
		utils.Clear()
		if retry == RETRY_WITH_DIFFERENT_CONFIG {
			interval = getInterval()
			vk = getVk()
		}
		retry = run(interval, vk)
	}
}
