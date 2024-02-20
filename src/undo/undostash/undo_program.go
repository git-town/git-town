package undostash

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

func DetermineUndoStashProgram(beginStashSize, endStashSize gitdomain.StashSize) program.Program {
	stashDiff := NewStashDiff(beginStashSize, endStashSize)
	return stashDiff.Program()
}
