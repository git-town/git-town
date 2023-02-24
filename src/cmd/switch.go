package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func switchCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "switch",
		Short: "Displays the local branches visually and allows switching between them",
		Run: func(cmd *cobra.Command, args []string) {
			roots := repo.Config.BranchAncestryRoots()
			for _, root := range roots {
				printBranch(0, root, repo)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
	}
}

func printBranch(depth int, branch string, repo *git.ProdRepo) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	fmt.Println(branch)
	children := repo.Silent.Config.ChildBranches(branch)
	for _, child := range children {
		printBranch(depth+1, child, repo)
	}
}
