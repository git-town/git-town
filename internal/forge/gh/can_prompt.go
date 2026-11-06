package gh

import (
	"os"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
)

// CanPrompt indicates whether gh may prompt the user for missing data.
// This is the case if Git Town runs interactively and gh is able to prompt.
// The latter mirrors IOStreams.CanPrompt in gh: STDIN and STDOUT must be terminals
// and GH_PROMPT_DISABLED must not be set (any value disables prompts).
// It ignores the "prompt" setting in the gh config.
type CanPrompt bool

func NewCanPrompt(interactive configdomain.Interactive, stdioIsTerminal configdomain.StdioIsTerminal) CanPrompt {
	if !interactive.IsEnabled() {
		return false
	}
	if _, disabled := os.LookupEnv("GH_PROMPT_DISABLED"); disabled {
		return false
	}
	return CanPrompt(stdioIsTerminal.IsTerminal())
}
