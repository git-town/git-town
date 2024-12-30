package print

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/colors"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/gohacks"
	"github.com/git-town/git-town/v17/internal/messages"
)

func Footer(verbose configdomain.Verbose, commandsCount gohacks.Counter, finalMessages []string) {
	fmt.Println()
	if verbose {
		fmt.Printf(messages.CommandsRun, commandsCount)
	}
	for _, finalMessage := range finalMessages {
		fmt.Println("\n" + colors.Cyan().Styled(finalMessage))
	}
}
