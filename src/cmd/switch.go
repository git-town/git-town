package cmd

import (
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func switchCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "switch",
		Short: "Displays the local branches visually and allows switching between them",
		Run: func(cmd *cobra.Command, args []string) {
			roots := repo.Config.BranchAncestryRoots()
			if err := keyboard.Open(); err != nil {
				cli.Exit(err)
			}
			cursor := uint8(0)
			for {
				for _, root := range roots {
					printBranch(printOptions{
						indent: 0,
						cursor: cursor,
						branch: root,
						repo:   repo,
					})
				}
				char, key, err := keyboard.GetSingleKey()
				if err != nil {
					cli.Exit(err)
				}
				switch char {
				case 'j':
					cursor += 1
				case 'k':
					cursor -= 1
				}
				if key == keyboard.KeyEsc {
					break
				}
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
	pos    uint8  // the current position in the list
	indent uint8  // the indentation of the current item
	cursor uint8  // position of the cursor in the list
	branch string // text of the list item
	repo   *git.ProdRepo
}

func printBranch(args printOptions) {
	space := "  "
	for i := uint8(0); i < args.indent; i++ {
		space += "  "
	}
	if args.cursor == args.pos {
		space = "*" + space[1:]
	}
	fmt.Println(space + args.branch)
	children := args.repo.Silent.Config.ChildBranches(args.branch)
	for _, child := range children {
		args.pos++
		printBranch(printOptions{
			pos:    args.pos,
			indent: args.indent + 1,
			cursor: args.cursor,
			branch: child,
			repo:   args.repo,
		})
	}
}
