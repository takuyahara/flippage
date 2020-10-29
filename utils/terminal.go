package utils

import (
	"fmt"
)

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func ClearLine() {
	fmt.Print("\033[2K\r")
}

func ClearPreviousLine() {
	fmt.Print("\033[0A\033[2K\r")
}
