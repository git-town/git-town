// Package runstate stores already executed and to be executed commands on disk.
// This is used by the `abort`, `continue`, and `undo` commands.
//
// When a Git Town command finishes, it stores the steps to undo or continue itself
// in a _state file_ located in the temp folder of the machine. The struct that
// holds the state is in [run_state.go](../src/runstate/run_state.go).
package runstate
