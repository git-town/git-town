package cmd

import (
	"sort"
	"strings"

	"atomicgo.dev/cursor"
	"github.com/eiannone/keyboard"
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
			input, err := createInput(currentBranch, 0, repo)
			if err != nil {
				cli.Exit(err)
			}
			cursor.Hide()
			defer func() {
				cursor.Show()
				keyboard.Close()
			}()
			for {
				input.Display()
				err := input.HandleInput()
				if err != nil {
					cli.Exit(err)
				}
				if input.Status == dialog.ModalInputStatusAborted {
					break
				}
			}
			input.Display()
			if input.Status == dialog.ModalInputStatusSelected {
				if input.SelectedValue() != currentBranch {
					err = repo.Silent.CheckoutBranch(input.SelectedValue())
					if err != nil {
						cli.Exit(err)
					}
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

func createInput(currentBranch string, indent int, repo *git.ProdRepo) (*dialog.ModalInput, error) {
	roots := repo.Config.BranchAncestryRoots()
	if err := keyboard.Open(); err != nil {
		return nil, err
	}
	entries := []dialog.ModalEntry{}
	var err error
	for _, root := range roots {
		entries, err = addEntries(entries, root, 0, repo)
		if err != nil {
			return nil, err
		}
	}
	cursorPos := 0
	for e, entry := range entries {
		if entry.Value == currentBranch {
			cursorPos = e
			break
		}
	}
	return &dialog.ModalInput{
		Entries:    entries,
		CursorPos:  uint8(cursorPos),
		CursorText: "> ",
	}, nil
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
