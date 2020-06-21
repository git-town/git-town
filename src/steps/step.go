package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// Step represents a dedicated activity within a Git Town command.
// Git Town commands are comprised of a number of steps that need to be executed.
type Step interface {
	CreateAbortStep() Step
	CreateContinueStep() Step
	CreateUndoStep(*git.ProdRepo) (Step, error)
	GetAutomaticAbortError() error
	Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error
	ShouldAutomaticallyAbortOnError() bool
}
