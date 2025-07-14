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

func enterGitLabToken() *cobra.Command {
	return &cobra.Command{
		Use: "gitlab-token",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.GitLabToken]{
				ConfigFileValue: None[forgedomain.GitLabToken](),
				HelpText:        dialog.GitLabTokenHelp,
				Inputs:          dialogInputs,
				LocalValue:      None[forgedomain.GitLabToken](),
				ParseFunc:       dialog.WrapParseFunc(forgedomain.ParseGitLabToken),
				Prompt:          "Your GitLab token: ",
				ResultMessage:   messages.GitLabToken,
				Title:           dialog.GitLabTokenTitle,
				UnscopedValue:   None[forgedomain.GitLabToken](),
			})
			return err
		},
	}
}
