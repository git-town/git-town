package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v23/internal/messages"
	"github.com/git-town/git-town/v23/internal/vm/shared"
)

type ExitToShellIfUncommittedChanges struct{}

func (self *ExitToShellIfUncommittedChanges) Run(args shared.RunArgs) error {
	uncommittedFiles, err := args.Git.UncommittedFiles(args.Backend)
	if err != nil {
		return err
	}
	if len(uncommittedFiles) > 0 {
		return ErrWalkUncommittedChanges
	}
	return nil
}

// ErrWalkUncommittedChanges indicates uncommitted changes in the repo
// that the user needs to commit before the Git Town command can continue.
var ErrWalkUncommittedChanges = errors.New(messages.WalkUncommittedChanges)
