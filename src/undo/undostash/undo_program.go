package undostash

import (
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/program"
)

func DetermineUndoStashProgram(initialStashSnapshot undodomain.StashSnapshot, backend *git.BackendCommands) (program.Program, error) {
	finalStashSnapshot, err := backend.StashSnapshot()
	if err != nil {
		return program.Program{}, err
	}
	stashDiff := NewStashDiff(initialStashSnapshot, finalStashSnapshot)
	return stashDiff.Program(), nil
}
