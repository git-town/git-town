package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
	"github.com/spf13/cobra"
)

var repoCommand = &cobra.Command{
	Use:   "repo",
	Short: "Opens the repository homepage",
	Long: `Opens the repository homepage

Supported for repositories hosted on GitHub, GitLab, and Bitbucket.
Derives the Git provider from the "origin" remote.
You can override this detection with
"git config git-town.code-hosting-driver <DRIVER>"
where DRIVER is "github", "gitlab", or "bitbucket".

When using SSH identities, run
"git config git-town.code-hosting-origin-hostname <HOSTNAME>"
where HOSTNAME matches what is in your ssh config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := git.NewProdRepo()
		driver := drivers.Get(repo.Configuration)
		if driver == nil {
			fmt.Println(drivers.UnsupportedHostingError())
			os.Exit(1)
		}
		script.OpenBrowser(driver.GetRepositoryURL())
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		if err := validateIsConfigured(); err != nil {
			return err
		}
		if err := git.Config().ValidateIsOnline(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(repoCommand)
}
