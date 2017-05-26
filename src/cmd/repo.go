package cmd

import (
	"github.com/Originate/git-town/src/drivers"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/spf13/cobra"
)

var repoCommand = &cobra.Command{
	Use:   "repo",
	Short: "Opens the repository homepage",
	Long: `Opens the repository homepage

Supported only for repositories hosted on GitHub, GitLab, and Bitbucket.
When using hosted versions of GitHub, GitLab, or Bitbucket,
make sure that your SSH identity contains the phrase "github", "gitlab", or
 "bitbucket", so that Git Town can guess which hosting service you use.

Example: your SSH identity should be something like
         "git@github-as-account1:Originate/git town.git"`,
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		prompt.EnsureIsConfigured()
		driver := drivers.GetCodeHostingDriver()
		repository := git.GetURLRepositoryName(git.GetRemoteOriginURL())
		script.OpenBrowser(driver.GetRepositoryURL(repository))
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func init() {
	RootCmd.AddCommand(repoCommand)
}
