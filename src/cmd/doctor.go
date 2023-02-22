package cmd

import (
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/spf13/cobra"
)

func doctorCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Displays diagnostic information",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Printf("Git Town v%s, built %s\n\n", version, buildDate)
			isRepo := repo.Silent.IsRepository()
			if isRepo {
				cli.Println("Git repository detected")
				// print origin
				// print upstream
			} else {
				cli.Println("The current folder does not seem to contain a Git repository")
			}
			cli.Print("- main branch: ")
			mainbranch := repo.Config.MainBranch()
			if mainbranch == "" {
				cli.Println("(not configured)")
			} else {
				cli.Println(mainbranch)
			}
			cli.Print("- perennial branches: ")
			perennialBranches := repo.Config.PerennialBranches()
			if len(perennialBranches) == 0 {
				cli.Println("(none)")
			} else {
				cli.Println(strings.Join(perennialBranches, ", "))
			}
			cli.Print("\nHosting service: ")
			driver, err := hosting.NewDriver(&repo.Config, &repo.Silent, cli.PrintDriverAction)
			if err == nil {
				cli.Println(driver.HostingServiceName())
				cli.Println("- repo URL: ", driver.RepositoryURL())
				cli.Println("- API token: %s", repo.)
			} else {
				cli.Println("(cannot determine: %v)", err)
			}
		},
		Args: cobra.NoArgs,
	}
}
