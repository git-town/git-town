package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterGitHubToken() *cobra.Command {
	return &cobra.Command{
		Use: "github-token",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.GitHubToken]{
				ConfigFileValue: None[forgedomain.GitHubToken](),
				HelpText:        dialog.GitHubTokenHelp,
				Inputs:          dialogInputs,
				LocalValue:      None[forgedomain.GitHubToken](),
				ParseFunc:       dialog.WrapParseFunc(forgedomain.ParseGitHubToken),
				Prompt:          "Your GitHub token: ",
				ResultMessage:   messages.GitHubToken,
				Title:           dialog.GitHubTokenTitle,
				UnscopedValue:   None[forgedomain.GitHubToken](),
			})
			return err
		},
	}
}
