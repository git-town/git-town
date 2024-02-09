package print

import (
	"fmt"

	"github.com/muesli/termenv"
)

// The Logger logger logs activities of a particular component on the CLI.
type Logger struct{}

func (l Logger) Failed(failure error) {
	boldRed := termenv.String().Bold().Foreground(termenv.ANSIRed)
	fmt.Println(boldRed.Styled(fmt.Sprintf("FAILED: %v\n", failure)))
}

func (l Logger) Start(template string, data ...interface{}) {
	fmt.Println()
	fmt.Println(Bold.Styled(fmt.Sprintf(template, data...)))
}

func (l Logger) Success() {
	boldGreen := termenv.String().Bold().Foreground(termenv.ANSIGreen)
	fmt.Println(boldGreen.Styled("ok\n"))
}
