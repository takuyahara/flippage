package utils

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/takuyahara/flippage/config"
)

var width int
var height int

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func ClearLine() {
	fmt.Print("\033[2K\r")
}

func ClearPreviousLine(rep ...int) {
	repeat := 1
	if len(rep) == 1 {
		if rep[0] < 0 {
			panic(`rep must be greater equal than 0.`)
		}
		repeat = rep[0]
	} else if len(rep) > 1 {
		panic(`Cannot passed more than 1 argument.`)
	}
	fmt.Print(strings.Repeat("\033[0A\033[2K\r", repeat))
}

func GoToPreviousLine(rep ...int) {
	repeat := 1
	if len(rep) == 1 {
		if rep[0] < 0 {
			panic(`rep must be greater equal than 0.`)
		}
		repeat = rep[0]
	} else if len(rep) > 1 {
		panic(`Cannot passed more than 1 argument.`)
	}
	fmt.Print(strings.Repeat("\033[0A\r", repeat))
}

func GoToNextLine(rep ...int) {
	repeat := 1
	if len(rep) == 1 {
		if rep[0] < 0 {
			panic(`rep must be greater equal than 0.`)
		}
		repeat = rep[0]
	} else if len(rep) > 1 {
		panic(`Cannot passed more than 1 argument.`)
	}
	fmt.Print(strings.Repeat("\033[0B\r", repeat))
}

func DefineTerminalSize() {
	cmd := exec.Command(`stty`, `size`)
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	size := strings.Split(strings.TrimSpace(string(output)), ` `)
	height, err = strconv.Atoi(size[0])
	if err != nil {
		panic(err)
	}
	width, err = strconv.Atoi(size[1])
	if err != nil {
		panic(err)
	}
}

func ShowProgressBarLeft(progress float64, msg ...string) {
	message := ``
	messageExists := false
	if len(msg) == 1 {
		message = msg[0]
		messageExists = true
	} else if len(msg) > 1 {
		panic(`Cannot passed more than 2 arguments.`)
	}
	Clear()
	repeatBar := int(math.Round(float64(width) * progress))
	repeatSpace := width - repeatBar
	if messageExists {
		GoToNextLine(2)
		fmt.Print(strings.Repeat(` `, repeatSpace) + "\033[46;1m" + strings.Repeat(` `, repeatBar) + "\033[0m")
		GoToPreviousLine(2)
		fmt.Print(message)
	} else {
		fmt.Print(strings.Repeat(` `, repeatSpace) + "\033[46;1m" + strings.Repeat(` `, repeatBar) + "\033[0m")
	}
}

func ShowProgressBarRight(progress float64, msg ...string) {
	message := ``
	messageExists := false
	if len(msg) == 1 {
		message = msg[0]
		messageExists = true
	} else if len(msg) > 1 {
		panic(`Cannot passed more than 2 arguments.`)
	}
	Clear()
	repeatBar := int(math.Round(float64(width) * progress))
	if messageExists {
		GoToNextLine(2)
		fmt.Print("\033[46;1m" + strings.Repeat(` `, repeatBar) + "\033[0m")
		GoToPreviousLine(2)
		fmt.Print(message)
	} else {
		fmt.Print("\033[46;1m" + strings.Repeat(` `, repeatBar) + "\033[0m")
	}
}

func ShowProgressBarDown(progress float64, msg ...string) {
	message := ``
	if len(msg) == 1 {
		message = msg[0]
	} else if len(msg) > 1 {
		panic(`Cannot passed more than 2 arguments.`)
	}
	Clear()
	messages := strings.Split(message, "\n")
	repeatBar := int(math.Round(float64(height) * progress))
	for i, line := range messages {
		firstChar := ``
		msgMargin := config.MARGIN_VERTICAL_PROGRESS_BAR
		if i < repeatBar {
			firstChar = "\033[46;1m \033[0m"
			msgMargin--
		}
		fmt.Printf("%s%s%s\n", firstChar, strings.Repeat(` `, msgMargin), line)
	}
	repeatBarRemaining := repeatBar - len(messages)
	if repeatBarRemaining > 0 {
		fmt.Print(strings.Repeat("\033[46;1m \033[0m\n", repeatBarRemaining-1) + "\033[46;1m \033[0m")
	}
	// Move cursor back to top
	if repeatBarRemaining > 0 {
		GoToPreviousLine(repeatBarRemaining - 1)
	}
	GoToPreviousLine(len(messages))
	// Move cursor to message's last line
	if len(messages) > 0 {
		GoToNextLine(len(messages) - 1)
		firstChar := ``
		msgMargin := config.MARGIN_VERTICAL_PROGRESS_BAR
		if len(messages)-1 < repeatBar {
			firstChar = "\033[46;1m \033[0m"
			msgMargin--
		}
		fmt.Printf("%s%s%s", firstChar, strings.Repeat(` `, msgMargin), messages[len(messages)-1])
	}
}

func ShowProgressBarScroll(progress float64, msg ...string) {
	ShowProgressBarDown(progress, msg...)
}
