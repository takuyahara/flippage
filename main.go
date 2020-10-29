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
	"strings"

	"github.com/micmonay/keybd_event"
)

const APP_NAME = `Flippage`
const (
	RETRY_WITH_SAME_CONFIG = iota
	RETRY_WITH_DIFFERENT_CONFIG
	NO_RETRY
)
const (
	MODE_FLIP = iota
	MODE_SCROLL
)

func getInterval() int {
	var interval int
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(`Type interval in second which is non-zero uint: `)
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

func getMode() (int, int) {
	var mode int
	var vk int
	scanner := bufio.NewScanner(os.Stdin)
	for {
		utils.ClearPreviousLine()
		fmt.Print("Specify a direction to flip (\033[4ml\033[24meft/\033[4mr\033[24might/\033[4md\033[24mown/\033[4ms\033[24mcroll): ")
		scanner.Scan()
		scanned := strings.ToLower(scanner.Text())
		if regexp.MustCompile(`^(?:left|right|down|scroll|l|r|d|s)$`).MatchString(scanned) {
			switch scanned {
			case `l`, `left`:
				mode = MODE_FLIP
				vk = keybd_event.VK_LEFT
			case `r`, `right`:
				mode = MODE_FLIP
				vk = keybd_event.VK_RIGHT
			case `d`, `down`:
				mode = MODE_FLIP
				vk = keybd_event.VK_DOWN
			case `s`, `scroll`:
				mode = MODE_SCROLL
				vk = keybd_event.VK_UP
			}
			break
		}
		utils.ClearLine()
	}
	return mode, vk
}

func getRetry() uint {
	var retry uint
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(strings.Join([]string{
			"Do you want to retry? (S/d/q)? ",
			" [S] Retry with \033[4ms\033[24mame config",
			" [d] Retry with \033[4md\033[24mifferent config",
			" [q] Don't retry and \033[4mq\033[24muit app",
		}, "\n"))
		utils.GoToPreviousLine(3)
		fmt.Print(`Do you want to retry? (S/d/q)? `)
		scanner.Scan()
		scanned := strings.ToLower(scanner.Text())
		if regexp.MustCompile(`^(?:S|d|q)?$`).MatchString(scanned) {
			switch scanned {
			case ``, `s`:
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

func runFlip(interval int, vk int) uint {
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

func runScroll(interval int) uint {
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
	go listener.Scroll(interval)
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
	mode, vk := getMode()
	var retry uint
	if mode == MODE_FLIP {
		retry = runFlip(interval, vk)
	} else {
		retry = runScroll(interval)
	}
	for retry != NO_RETRY {
		utils.Clear()
		if retry == RETRY_WITH_DIFFERENT_CONFIG {
			interval = getInterval()
			mode, vk = getMode()
		}
		if mode == MODE_FLIP {
			retry = runFlip(interval, vk)
		} else {
			retry = runScroll(interval)
		}
	}
}
