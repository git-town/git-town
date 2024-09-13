package print

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/colors"
)

// The Logger logger logs activities of a particular component on the CLI.
type Logger struct{}

func (l Logger) Failed(failure error) {
	l.Log(colors.BoldRed().Styled(fmt.Sprintf("FAILED: %v\n", failure)))
}

func (l Logger) Log(text string) {
	fmt.Println(text)
}

func (l Logger) Ok() {
	l.Success("ok")
}

func (l Logger) Start(template string, data ...interface{}) {
	fmt.Println()
	fmt.Print(colors.Bold().Styled(fmt.Sprintf(template, data...)))
}

func (l Logger) Success(message string) {
	l.Log(colors.BoldGreen().Styled(message))
}
