package undostash

import (
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

func DetermineUndoStashProgram(initialStashSnapshot gitdomain.StashSize, backend *git.BackendCommands) (program.Program, error) {
	finalStashSnapshot, err := backend.StashSize()
	if err != nil {
		return program.Program{}, err
	}
	stashDiff := NewStashDiff(initialStashSnapshot, finalStashSnapshot)
	return stashDiff.Program(), nil
}
