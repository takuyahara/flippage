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
const (
	RETRY_WITH_SAME_CONFIG = iota
	RETRY_WITH_DIFFERENT_CONFIG
	NO_RETRY
)

func getInterval() int {
	var interval int
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
			case `l`:
			case `left`:
				vk = keybd_event.VK_LEFT
			case `r`:
			case `right`:
				vk = keybd_event.VK_RIGHT
			case `u`:
			case `up`:
				vk = keybd_event.VK_UP
			case `d`:
			case `down`:
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
		fmt.Print(`Do you want to retry? (S/d/n)? 
 [S] Retry with same config
 [d] Retry with different config
 [n] Don't retry and close app`)
		utils.GoToPreviousLine(3)
		fmt.Print(`Do you want to retry? (S/d/n)? `)
		scanner.Scan()
		scanned := scanner.Text()
		if regexp.MustCompile(`^(?:S|d|n)?$`).MatchString(scanned) {
			switch scanned {
			case ``:
			case `S`:
				retry = RETRY_WITH_SAME_CONFIG
			case `d`:
				retry = RETRY_WITH_DIFFERENT_CONFIG
			case `n`:
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

func run(interval int, vk int) uint {
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
	keyboard.Stop()
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
