package cmd

import (
	"sort"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func switchCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "switch",
		Short: "Displays the local branches visually and allows switching between them",
		Run: func(cmd *cobra.Command, args []string) {
			currentBranch, err := repo.Silent.CurrentBranch()
			if err != nil {
				cli.Exit(err)
			}
			input, cleanup, err := createInput(currentBranch, 0, repo)
			if err != nil {
				cli.Exit(err)
			}
			defer cleanup()
			userChoice, err := input.Display()
			if err != nil {
				cli.Exit(err)
			}
			if userChoice != nil && *userChoice != currentBranch {
				err = repo.Silent.CheckoutBranch(*userChoice)
				if err != nil {
					cli.Exit(err)
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

func createInput(currentBranch string, indent int, repo *git.ProdRepo) (*dialog.ModalInput, func(), error) {
	roots := repo.Config.BranchAncestryRoots()
	entries := []dialog.ModalEntry{}
	var err error
	for _, root := range roots {
		entries, err = addEntries(entries, root, 0, repo)
		if err != nil {
			return nil, func() {}, err
		}
	}
	return dialog.NewModalInput(entries, "> ", currentBranch)
}

func addEntries(entries []dialog.ModalEntry, branch string, indent int, repo *git.ProdRepo) ([]dialog.ModalEntry, error) {
	entries = append(entries, dialog.ModalEntry{
		Text:  strings.Repeat("  ", indent) + branch,
		Value: branch,
	})
	children := repo.Silent.Config.ChildBranches(branch)
	sort.Strings(children)
	var err error
	for _, child := range children {
		entries, err = addEntries(entries, child, indent+1, repo)
		if err != nil {
			return entries, err
		}
	}
	return entries, nil
}
