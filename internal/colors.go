package internal

import (
	"fmt"

	"github.com/fatih/color"
)

func InitColors() {
	color.NoColor = false
}

func Red(s string) string    { return color.New(color.FgRed).Sprint(s) }
func Green(s string) string  { return color.New(color.FgGreen).Sprint(s) }
func Yellow(s string) string { return color.New(color.FgYellow).Sprint(s) }
func Blue(s string) string   { return color.New(color.FgCyan).Sprint(s) }
func Gray(s string) string   { return color.New(color.FgHiBlack).Sprint(s) }
func White(s string) string  { return color.New(color.FgWhite).Sprint(s) }
func Bold(s string) string   { return color.New(color.Bold).Sprint(s) }

func Separator() string {
	return Gray("─────────────────────────────────────")
}

func Success(s string) {
	fmt.Println(Green("✓ " + s))
}

func Error(s string) {
	fmt.Println(Red("✗ " + s))
}
