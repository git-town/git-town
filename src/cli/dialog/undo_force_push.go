package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

const (
	undoForcePushTitle = `force-push to remote branch`
	undoForcePushHelp  = `
Should I force-push remote branch %q?

Existing commit: %q
Commit to be pushed: %q

`
)

type YesNoEntry struct {
	Text  string
	Value bool
}

// GitHubToken lets the user enter the GitHub API token.
func ForcePushBranch(branch gitdomain.RemoteBranchName, oldValue, newValue string, inputs components.TestInput) (Option[configdomain.GitHubToken], bool, error) {
	text, aborted, err := components.RadioList(list.Entries[string]{})(components.TextFieldArgs{
		ExistingValue: oldValue,
		Help:          gitHubTokenHelp,
		Prompt:        "Your GitHub API token: ",
		TestInput:     inputs,
		Title:         githubTokenTitle,
	})
	fmt.Printf(messages.GitHubToken, components.FormattedSecret(text, aborted))
	return configdomain.NewGitHubTokenOption(text), aborted, err
}
