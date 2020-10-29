package utils

import (
	"fmt"
	"strings"
)

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func ClearLine() {
	fmt.Print("\033[2K\r")
}

func ClearPreviousLine(rep ...int) {
	repeat := 1
	if len(rep) == 1 {
		if rep[0] <= 0 {
			panic(`rep must be greater than 0.`)
		}
		repeat = rep[0]
	} else if len(rep) > 1 {
		panic(`Cannot passed more than 1 argument.`)
	}
	fmt.Printf(strings.Repeat("\033[0A\033[2K\r", repeat))
}

func GoToPreviousLine(rep ...int) {
	repeat := 1
	if len(rep) == 1 {
		if rep[0] <= 0 {
			panic(`rep must be greater than 0.`)
		}
		repeat = rep[0]
	} else if len(rep) > 1 {
		panic(`Cannot passed more than 1 argument.`)
	}
	fmt.Printf(strings.Repeat("\033[0A\r", repeat))
}

func GoToNextLine(rep ...int) {
	repeat := 1
	if len(rep) == 1 {
		if rep[0] <= 0 {
			panic(`rep must be greater than 0.`)
		}
		repeat = rep[0]
	} else if len(rep) > 1 {
		panic(`Cannot passed more than 1 argument.`)
	}
	fmt.Printf(strings.Repeat("\033[0B\r", repeat))
}
