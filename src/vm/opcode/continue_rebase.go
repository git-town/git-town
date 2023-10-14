package opcode

// ContinueRebase finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebase struct {
	undeclaredOpcodeMethods
}

func (step *ContinueRebase) CreateAbortProgram() []Opcode {
	return []Opcode{
		&AbortRebase{},
	}
}

func (step *ContinueRebase) CreateContinueProgram() []Opcode {
	return []Opcode{
		step,
	}
}

func (step *ContinueRebase) Run(args RunArgs) error {
	repoStatus, err := args.Runner.Backend.RepoStatus()
	if err != nil {
		return err
	}
	if repoStatus.RebaseInProgress {
		return args.Runner.Frontend.ContinueRebase()
	}
	return nil
}
