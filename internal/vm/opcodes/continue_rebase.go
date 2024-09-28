package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

// ContinueRebaseIfNeeded finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseIfNeeded struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ContinueRebaseIfNeeded) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortRebase{},
	}
}

func (self *ContinueRebaseIfNeeded) Run(args shared.RunArgs) error {
	repoStatus, err := args.Git.RepoStatus(args.Backend)
	if err != nil {
		return err
	}
	if repoStatus.RebaseInProgress {
		return args.Git.ContinueRebase(args.Frontend)
	}
	return nil
}
