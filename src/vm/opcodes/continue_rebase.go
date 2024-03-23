package opcodes

import "github.com/git-town/git-town/v13/src/vm/shared"

// ContinueRebase finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebase struct {
	undeclaredOpcodeMethods
}

func (self *ContinueRebase) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortRebase{},
	}
}

func (self *ContinueRebase) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *ContinueRebase) Run(args shared.RunArgs) error {
	repoStatus, err := args.Runner.Backend.RepoStatus()
	if err != nil {
		return err
	}
	if repoStatus.RebaseInProgress {
		return args.Runner.Frontend.ContinueRebase()
	}
	return nil
}
