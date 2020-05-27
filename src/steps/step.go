package steps

import "github.com/git-town/git-town/src/git"

// Step represents a dedicated activity within a Git Town command.
// Git Town commands are comprised of a number of steps that need to be executed.
type Step interface {
	CreateAbortStep() Step
	CreateContinueStep() Step
	CreateUndoStep() Step
	GetAutomaticAbortErrorMessage() string
	Run(repo *git.ProdRepo) error
	ShouldAutomaticallyAbortOnError() bool
}
