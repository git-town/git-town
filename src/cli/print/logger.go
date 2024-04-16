package print

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/colors"
)

// The Logger logger logs activities of a particular component on the CLI.
type Logger struct{}

func (l Logger) Failed(failure error) {
	fmt.Println(colors.BoldRed().Styled(fmt.Sprintf("FAILED: %v\n", failure)))
}

func (l Logger) Start(template string, data ...interface{}) {
	fmt.Println()
	fmt.Print(colors.Bold().Styled(fmt.Sprintf(template, data...)))
}

func (l Logger) Success() {
	fmt.Println(colors.BoldGreen().Styled("ok"))
}
