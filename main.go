package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/takuyahara/flippage/appinfo"
	"github.com/takuyahara/flippage/config"
	"github.com/takuyahara/flippage/listener"
	"github.com/takuyahara/flippage/utils"
)

const APP_NAME = `Flippage`

func getInterval() int {
	var interval int
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\033[1mType interval in second which is non-zero uint: \033[m")
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

func getConfig() (int, int, string) {
	var direction int
	var mode int
	var key string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		utils.ClearPreviousLine()
		fmt.Print("\033[1mSpecify a direction to flip (\033[4ml\033[24meft/\033[4mr\033[24might/\033[4md\033[24mown/\033[4ms\033[24mcroll): \033[m")
		scanner.Scan()
		scanned := strings.ToLower(scanner.Text())
		if regexp.MustCompile(`^(?:left|right|down|scroll|l|r|d|s)$`).MatchString(scanned) {
			switch scanned {
			case `l`, `left`:
				direction = config.DIRECTION_LEFT
				mode = config.MODE_FLIP
				key = "left"
			case `r`, `right`:
				direction = config.DIRECTION_RIGHT
				mode = config.MODE_FLIP
				key = "right"
			case `d`, `down`:
				direction = config.DIRECTION_DOWN
				mode = config.MODE_FLIP
				key = "down"
			case `s`, `scroll`:
				direction = config.DIRECTION_DOWN
				mode = config.MODE_SCROLL
				key = "up"
			}
			break
		}
		utils.ClearLine()
	}
	return direction, mode, key
}

func getRetry() uint {
	var retry uint
	scanner := bufio.NewScanner(os.Stdin)
	for {
		msgPrompt := []string{
			"\033[1mDo you want to retry? (S/d/q)? \033[m",
			" [S] Retry with \033[4ms\033[24mame config",
			" [d] Retry with \033[4md\033[24mifferent config",
			" [q] Don't retry and \033[4mq\033[24muit app",
		}
		fmt.Print(strings.Join(msgPrompt, "\n"))
		utils.GoToPreviousLine(3)
		fmt.Print(msgPrompt[0])
		scanner.Scan()
		scanned := strings.ToLower(scanner.Text())
		if regexp.MustCompile(`^(?:s|d|q)?$`).MatchString(scanned) {
			switch scanned {
			case ``, `s`:
				retry = config.RETRY_WITH_SAME_CONFIG
			case `d`:
				retry = config.RETRY_WITH_DIFFERENT_CONFIG
			case `q`:
				retry = config.NO_RETRY
			}
			break
		}
		utils.GoToNextLine(3)
		utils.ClearPreviousLine(4)
	}
	utils.GoToNextLine(3)
	return retry
}

func runFlip(direction, interval int, key string) uint {
	timeStart := time.Now()
	chAppInfoForeground := make(chan appinfo.AppInfo, 2)
	// Wait until foreground app changes
	utils.ClearPreviousLine()
	fmt.Print(`Waiting to foreground app changes...`)
	appinfo.GetChangedForegroundInfo(chAppInfoForeground)
	appInfoTarget := <-chAppInfoForeground
	// Flip page automatically
	utils.ClearLine()
	utils.DefineTerminalSize()
	message := fmt.Sprintf("%s has activated for %s.\n", APP_NAME, appInfoTarget.Name)
	listener.ListenEvents()
	go listener.Flip(message, direction, interval, key)
	appinfo.GetChangedForegroundInfo(chAppInfoForeground)
	<-chAppInfoForeground
	// Close app
	timeEnd := time.Now()
	timeSpent := formatElapsedDuration(timeEnd.Sub(timeStart))
	listener.Stop()
	utils.Clear()
	fmt.Print(message)
	fmt.Printf("%s has exited after spending %s as foreground app has changed.\n", APP_NAME, timeSpent)
	return getRetry()
}

func runScroll(interval int) uint {
	timeStart := time.Now()
	chAppInfoForeground := make(chan appinfo.AppInfo, 2)
	// Wait until foreground app changes
	utils.ClearPreviousLine()
	fmt.Print(`Waiting to foreground app changes...`)
	appinfo.GetChangedForegroundInfo(chAppInfoForeground)
	appInfoTarget := <-chAppInfoForeground
	// Flip page automatically
	utils.ClearLine()
	utils.DefineTerminalSize()
	message := fmt.Sprintf("%s has activated for %s.\n", APP_NAME, appInfoTarget.Name)
	listener.ListenEvents()
	go listener.Scroll(message, interval)
	appinfo.GetChangedForegroundInfo(chAppInfoForeground)
	<-chAppInfoForeground
	// Close app
	timeEnd := time.Now()
	timeSpent := formatElapsedDuration(timeEnd.Sub(timeStart))
	listener.Stop()
	utils.Clear()
	fmt.Print(message)
	fmt.Printf("%s has exited after spending %s as foreground app has changed.\n", APP_NAME, timeSpent)
	return getRetry()
}

func formatElapsedDuration(d time.Duration) string {
	elapsed := d.Truncate(time.Second)
	h := elapsed / time.Hour
	elapsed -= h * time.Hour
	m := elapsed / time.Minute
	elapsed -= m * time.Minute
	s := elapsed / time.Second
	dFormatted := fmt.Sprintf(`%dh%dm%ds`, h, m, s)
	if h == 0 {
		if m > 0 {
			dFormatted = fmt.Sprintf(`%dm%ds`, m, s)
		} else {
			dFormatted = fmt.Sprintf(`%ds`, s)
		}
	}
	return dFormatted
}

func main() {
	utils.Clear()
	interval := getInterval()
	direction, mode, key := getConfig()
	var retry uint
	if mode == config.MODE_FLIP {
		retry = runFlip(direction, interval, key)
	} else {
		retry = runScroll(interval)
	}
	for retry != config.NO_RETRY {
		utils.Clear()
		if retry == config.RETRY_WITH_DIFFERENT_CONFIG {
			interval = getInterval()
			direction, mode, key = getConfig()
		}
		if mode == config.MODE_FLIP {
			retry = runFlip(direction, interval, key)
		} else {
			retry = runScroll(interval)
		}
	}
}
