package steps

import (
	"log"
	"os"

	"github.com/Originate/git-town/lib/git"
)

// WrapOptions represents the options given to Wrap.
type WrapOptions struct {
	RunInGitRoot     bool
	StashOpenChanges bool
}

// Wrap wraps the given StepList in steps that
// change to the Git root directory or stash away open changes.
func Wrap(stepList StepList, options WrapOptions) (result StepList) {
	result.AppendList(stepList)

	if options.StashOpenChanges && git.HasOpenChanges() {
		result.Prepend(StashOpenChangesStep{})
		result.Append(RestoreOpenChangesStep{})
	}

	// TODO echo "preserve_checkout_history $INITIAL_PREVIOUS_BRANCH_NAME $INITIAL_BRANCH_NAME"

	initialDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	gitRootDirectory := git.GetRootDirectory()

	if options.RunInGitRoot && initialDirectory != gitRootDirectory {
		result.Prepend(ChangeDirectoryStep{Directory: gitRootDirectory})
		result.Append(ChangeDirectoryStep{Directory: initialDirectory})
	}

	return
}
