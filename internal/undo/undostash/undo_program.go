package undostash

import (
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/program"
)

func DetermineUndoStashProgram(beginStashSize, endStashSize gitdomain.StashSize) program.Program {
	return NewStashDiff(beginStashSize, endStashSize).Program()
}
