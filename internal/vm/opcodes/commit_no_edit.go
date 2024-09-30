package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type CommitNoEdit struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitNoEdit) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, None[gitdomain.CommitMessage](), true, None[gitdomain.Author]())
}
