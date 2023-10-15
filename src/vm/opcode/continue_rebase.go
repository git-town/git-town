package opcode

import "github.com/git-town/git-town/v9/src/vm/shared"

// ContinueRebase finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebase struct {
	undeclaredOpcodeMethods
}

func (op *ContinueRebase) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortRebase{},
	}
}

func (op *ContinueRebase) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		op,
	}
}

func (op *ContinueRebase) Run(args shared.RunArgs) error {
	repoStatus, err := args.Runner.Backend.RepoStatus()
	if err != nil {
		return err
	}
	if repoStatus.RebaseInProgress {
		return args.Runner.Frontend.ContinueRebase()
	}
	return nil
}
