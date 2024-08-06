package print

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/colors"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/gohacks"
	"github.com/git-town/git-town/v15/internal/messages"
)

func Footer(verbose configdomain.Verbose, commandsCount gohacks.Counter, finalMessages []string) {
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
