// Package runstate stores to be executed as well as already executed steps on disk.
// This is used by the "abort", "continue", and "undo" commands.
//
// Git Town represents individual steps within a Git Town command via the command pattern
// (https://en.wikipedia.org/wiki/Command_pattern). This allows fully automated and robust
// implementations of the "continue", "abort", and "undo" commands.
// When a Git Town command finishes, it stores the steps to undo or continue itself
// in a state file located in the temp folder of the machine.
// The logic for this is in run_state_to_disk.go.
//
package runstate
