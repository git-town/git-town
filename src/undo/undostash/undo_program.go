package undostash

import (
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

func DetermineUndoStashProgram(initialStashSize gitdomain.StashSize, backend *git.BackendCommands) (program.Program, error) {
	finalStashSize, err := backend.StashSize()
	if err != nil {
		return program.Program{}, err
	}
	stashDiff := NewStashDiff(initialStashSize, finalStashSize)
	return stashDiff.Program(), nil
}
