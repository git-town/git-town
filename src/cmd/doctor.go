package cmd

import (
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
			cli.Println("Git Town:")
			cli.Printf("- version: %s\n", version)
			cli.Printf("- built: %s\n\n", buildDate)
			isRepo := repo.Silent.IsRepository()
			if isRepo {
				cli.Println("Git repository:")
				// print origin
				// print upstream
			} else {
				cli.Println("The current folder does not seem to contain a Git repository")
			}
			cli.Print("- main branch: ")
			isOffline, err := repo.Config.IsOffline()
			if err == nil {
				if isOffline {
					cli.Println("enabled")
				} else {
					cli.Println("disabled")
				}
			} else {
				cli.Println("(cannot determine: %v)", err)
			}
			cli.Print("- main branch: ")
			mainbranch := repo.Config.MainBranch()
			if mainbranch == "" {
				cli.Println("(not configured)")
			} else {
				cli.Println(mainbranch)
			}
			cli.Println("\nhosting service:")
			driver, err := hosting.NewDriver(&repo.Config, &repo.Silent, cli.PrintDriverAction)
			if err == nil {
				cli.Println("- name: ", driver.HostingServiceName())
				cli.Println("- repo: ", driver.RepositoryURL())
			} else {
				cli.Println("(cannot determine: %v)", err)
			}
		},
		Args: cobra.NoArgs,
	}
}
