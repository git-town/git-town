package print

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/messages"
)

func Footer(verbose bool, commandsCount int, finalMessages []string) {
	fmt.Println()
	if verbose {
		fmt.Printf(messages.CommandsRun, commandsCount)
	}
	Messages(finalMessages)
}

// NoFinalMessages can be used by callers of PrintFooter to indicate
// that the command has no final messages to print.
var NoFinalMessages = []string{} //nolint:gochecknoglobals

// Messages prints the given messages to the user.
func Messages(messages []string) {
	for _, message := range messages {
		fmt.Println("\n" + colors.Cyan().Styled(message))
	}
}
