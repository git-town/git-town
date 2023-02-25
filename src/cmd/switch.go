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
				printBranch(printOptions{
					depth:  0,
					cursor: 0,
					branch: root,
					repo:   repo,
				})
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

type printOptions struct {
	depth  uint8
	cursor uint8
	branch string
	repo   *git.ProdRepo
}

func printBranch(args printOptions) {
	for i := uint8(0); i < args.depth; i++ {
		fmt.Print("  ")
	}
	fmt.Println(args.branch)
	children := args.repo.Silent.Config.ChildBranches(args.branch)
	for _, child := range children {
		printBranch(printOptions{
			depth:  args.depth + 1,
			branch: child,
			cursor: args.cursor,
			repo:   args.repo})
	}
}
