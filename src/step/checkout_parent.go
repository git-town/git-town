package step

// CheckoutParent checks out the parent of the current branch.
type CheckoutParent struct {
	Empty
}

func (step *CheckoutParent) Run(args RunArgs) error {
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	parent := args.Lineage.Parent(currentBranch)
	if currentBranch == parent {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(parent)
}
