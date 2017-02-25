package step

import (
  "log"
  "os"

  "github.com/Originate/gt/cmd/git"
)

type WrapOptions struct {
  RunInGitRoot bool
  StashOpenChanges bool
}

func Wrap(steps []Step, options WrapOptions) []Step {
  if options.StashOpenChanges && git.HasOpenChanges() {
    steps = append([]Step{new(StashOpenChangesStep)}, steps...)
    steps = append(steps, new(RestoreOpenChangesStep))
  }

  // TODO echo "preserve_checkout_history $INITIAL_PREVIOUS_BRANCH_NAME $INITIAL_BRANCH_NAME"

  initialDirectory, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }
  gitRootDirectory := git.GetRootDirectory()

  if options.RunInGitRoot && initialDirectory != gitRootDirectory {
    steps = append([]Step{ChangeDirectoryStep{directory: gitRootDirectory}}, steps...)
    steps = append(steps, ChangeDirectoryStep{directory: initialDirectory})
  }

  return steps
}
