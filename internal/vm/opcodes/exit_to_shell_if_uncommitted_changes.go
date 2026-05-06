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

// WalkUncommittedChangesError is a sentinel error
// that indicates that the "walk" command has detected uncommitted changes
// and wants to exit to the shell to allow the user to commit them
// before continuing to walk the next branch.
var WalkUncommittedChangesError = errors.New(messages.WalkUncommittedChanges)
