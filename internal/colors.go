package internal

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Red    = color.New(color.FgRed).SprintFunc()
	Green  = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Blue   = color.New(color.FgCyan).SprintFunc()
	Gray   = color.New(color.FgHiBlack).SprintFunc()
	White  = color.New(color.FgWhite).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
)

func InitColors() {} // больше не нужна, fatih/color делает всё сам

func Separator() string {
	return Gray("─────────────────────────────────────")
}

func Success(s string) {
	fmt.Println(Green("✓ " + s))
}

func Error(s string) {
	fmt.Println(Red("✗ " + s))
}
