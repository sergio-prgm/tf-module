package util

import "fmt"

type Color int

// TODO search ANSI bold, underline, background, etc.
const (
	Reset Color = iota
	Red
	Green
	Yellow
	Blue
	Purple
	Cyan
	Gray
	White
)

var clrs = map[Color]string{
	Reset:  "0m",
	Red:    "31m",
	Green:  "32m",
	Yellow: "33m",
	Blue:   "34m",
	Purple: "35m",
	Cyan:   "36m",
	Gray:   "37m",
	White:  "97m",
}

type Mode int

const (
	Normal Mode = iota
	Bold
	Underline
)

func clrBuilder(clr Color, mod Mode) string {
	return fmt.Sprintf("\033[%d;%s", mod, clrs[clr])
}

func EmphasizeStr(str string, clr Color, mod Mode) string {
	reset := clrBuilder(Reset, Normal)
	color := clrBuilder(clr, mod)
	return fmt.Sprintf("%s%s%s", color, str, reset)
}
