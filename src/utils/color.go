package utils

import "fmt"

const (
	NONE    = 37
	BLACK   = 30
	RED     = 31
	GREEN   = 32
	YELLOW  = 33
	BLUE    = 34
	MAGENTA = 35
	CYAN    = 36
	WHITE   = 37
)

func Format(c int, text string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, text)
}
