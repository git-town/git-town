package dialog

import (
	"strings"

	"github.com/git-town/git-town/v11/src/config"
	"github.com/git-town/git-town/v11/src/domain"
)

// queryBranch lets the user select a new branch via a visual dialog.
// Indicates via `validSelection` whether the user made a valid selection.
func QueryBranch(currentBranch domain.LocalBranchName, lineage config.Lineage) (selection domain.LocalBranchName, validSelection bool, err error) {
	entries, err := createEntries(lineage, currentBranch)
	if err != nil {
		return domain.EmptyLocalBranchName(), false, err
	}
	choice, err := ModalSelect(entries, currentBranch.String())
	if err != nil {
		return domain.EmptyLocalBranchName(), false, err
	}
	if choice == nil {
		return domain.EmptyLocalBranchName(), false, nil
	}
	return domain.NewLocalBranchName(*choice), true, nil
}

// AddEntryAndChildren adds the given branch and all its child branches to the given entries collection.
func AddEntryAndChildren(entries ModalSelectEntries, branch domain.LocalBranchName, indent int, lineage config.Lineage) (ModalSelectEntries, error) {
	entries = append(entries, ModalSelectEntry{
		Text:  strings.Repeat("  ", indent) + branch.String(),
		Value: branch.String(),
	})
	var err error
	for _, child := range lineage.Children(branch) {
		entries, err = AddEntryAndChildren(entries, child, indent+1, lineage)
		if err != nil {
			return entries, err
		}
	}
	return entries, nil
}

// createEntries provides all the entries for the branch dialog.
func createEntries(lineage config.Lineage, currentBranch domain.LocalBranchName) (ModalSelectEntries, error) {
	entries := ModalSelectEntries{}
	var err error
	for _, root := range lineage.Roots() {
		entries, err = AddEntryAndChildren(entries, root, 0, lineage)
		if err != nil {
			return nil, err
		}
	}
	if len(entries) == 0 {
		entries = append(entries, ModalSelectEntry{
			Text:  string(currentBranch),
			Value: string(currentBranch),
		})
	}
	return entries, nil
}
