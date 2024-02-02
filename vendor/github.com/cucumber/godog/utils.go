package godog

import (
	"strings"
	"time"

	"github.com/cucumber/godog/colors"
)

var (
	red    = colors.Red
	redb   = colors.Bold(colors.Red)
	green  = colors.Green
	blackb = colors.Bold(colors.Black)
	yellow = colors.Yellow
	cyan   = colors.Cyan
	cyanb  = colors.Bold(colors.Cyan)
	whiteb = colors.Bold(colors.White)
)

// repeats a space n times
func s(n int) string {
	if n < 0 {
		n = 1
	}
	return strings.Repeat(" ", n)
}

var timeNowFunc = func() time.Time {
	return time.Now()
}

func trimAllLines(s string) string {
	var lines []string
	for _, ln := range strings.Split(strings.TrimSpace(s), "\n") {
		lines = append(lines, strings.TrimSpace(ln))
	}
	return strings.Join(lines, "\n")
}
