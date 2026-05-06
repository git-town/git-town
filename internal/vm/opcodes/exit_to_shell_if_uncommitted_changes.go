package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type ExitToShellIfUncommittedChanges struct{}

func (self *ExitToShellIfUncommittedChanges) Run(args shared.RunArgs) error {
	uncommittedFiles, err := args.Git.UncommittedFiles(args.Backend)
	if err != nil {
		return err
	}
	if len(uncommittedFiles) > 0 {
		return WalkUncommittedChangesError
	}
	return nil
}

// WalkUncommittedChangesError indicates uncommitted changes in the repo
// that the user needs to commit before the Git Town command can continue.
var WalkUncommittedChangesError = errors.New(messages.WalkUncommittedChanges)
