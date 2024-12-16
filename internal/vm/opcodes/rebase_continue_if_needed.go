package opcodes

import "github.com/git-town/git-town/v17/internal/vm/shared"

// RebaseContinueIfNeeded finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type RebaseContinueIfNeeded struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseContinueIfNeeded) Run(args shared.RunArgs) error {
	repoStatus, err := args.Git.RepoStatus(args.Backend)
	if err != nil {
		return err
	}
	if repoStatus.RebaseInProgress {
		args.PrependOpcodes(&RebaseContinue{})
	}
	return nil
}
