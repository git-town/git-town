// Package cmd defines the Git Town commands.
//
// Each Git Town command begins by inspecting the current state of the Git
// repository (which branch you are on, whether you have open changes). If there
// are no errors, it generates a StepList instance containing the steps to run.
//
// Steps, located in src/steps, implement the individual steps that
// each Git Town command performs. Examples are steps to
// change to a different Git branch or to pull updates for the current branch.
//
// When executing a step, Git Town asks the step to provide the undo step for it
// and appends that undo step to the "undo list" for the current Git Town command.
// If a Git command fails (for example due to a merge conflict), then the program
// saves the current runstate (the steps to abort and continue the current command)
// to disk, informs the user, and exits.
//
// When running "git town continue", Git Town loads the runstate and executes the list of remaining steps.
// When running "git town abort", Git Town loads the runstate and executes the list of abort steps.
// When running "git town undo", Git Town loads the runstate and executes its undo list.
package cmd
