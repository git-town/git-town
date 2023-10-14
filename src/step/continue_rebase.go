package step

// ContinueRebase finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebase struct {
	Empty
}

func (step *ContinueRebase) CreateAbortProgram() []Step {
	return []Step{
		&AbortRebase{},
	}
}

func (step *ContinueRebase) CreateContinueProgram() []Step {
	return []Step{
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
