package steps

import (
  "log"
  "os"

  "github.com/Originate/gt/lib/git"
)

type WrapOptions struct {
  RunInGitRoot bool
  StashOpenChanges bool
}

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
