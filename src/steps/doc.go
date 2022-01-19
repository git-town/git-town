// Package steps defines the individual CLI operations that Git Town can execute.
// All steps implement the Step interface defined in step.go.
// Git Town doesn't execute steps directly.
// It organizes all activities it wants to perform in a StepList and executes that StepList instance.
package steps
