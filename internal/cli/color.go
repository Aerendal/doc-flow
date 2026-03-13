package cli

import "os"

// ANSI color codes — no external dependencies.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// colorEnabled returns true when stdout is a terminal and NO_COLOR is not set.
func colorEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func colorize(color, s string) string {
	if !colorEnabled() {
		return s
	}
	return color + s + colorReset
}

func green(s string) string  { return colorize(colorGreen, s) }
func red(s string) string    { return colorize(colorRed, s) }
func yellow(s string) string { return colorize(colorYellow, s) }
func cyan(s string) string   { return colorize(colorCyan, s) }
func bold(s string) string   { return colorize(colorBold, s) }

func passLabel() string { return green("✓ PASS") }
func failLabel() string { return red("✗ FAIL") }
func warnLabel() string { return yellow("⚠ WARN") }
