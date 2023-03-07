package runstate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/config"
)

// ErrorChecker helps avoid excessive error checking
// while gathering a larger number of values through fallible operations.
//
// This is based on ideas outlined in https://go.dev/blog/errors-are-values.
// Please be aware that this an experimental idea.
// Using this technique can lead to executing logic that would normally not run,
// using potentially invalid data, and potentially leading to unexpected runtime exceptions and side effects.
// Use with care and only if it's abundantly clear and obvious that there are no negative side effects.
type ErrorChecker struct {
	Err error `exhaustruct:"optional"`
}

// Check registers the given error and indicates
// whether this ErrorChecker contains an error now.
func (ec *ErrorChecker) Check(err error) bool {
	if ec.Err == nil {
		ec.Err = err
	}
	return ec.Err != nil
}

// Fail registers the error constructed using the given format arguments.
func (ec *ErrorChecker) Fail(format string, a ...any) {
	ec.Check(fmt.Errorf(format, a...))
}

// Bool provides the bool part of the given fallible function result
// while registering the given error.
func (ec *ErrorChecker) Bool(value bool, err error) bool {
	ec.Check(err)
	return value
}

// HostingService provides the config.HostingService part of the given fallible function result
// while registering the given error.
func (ec *ErrorChecker) HostingService(value config.HostingService, err error) config.HostingService {
	ec.Check(err)
	return value
}

// PullBranchStrategy provides the string part of the given fallible function result
// while registering the given error.
func (ec *ErrorChecker) PullBranchStrategy(value config.PullBranchStrategy, err error) config.PullBranchStrategy {
	ec.Check(err)
	return value
}

// String provides the string part of the given fallible function result
// while registering the given error.
func (ec *ErrorChecker) String(value string, err error) string {
	ec.Check(err)
	return value
}

// Strings provides the []string part of the given fallible function result
// while registering the given error.
func (ec *ErrorChecker) Strings(value []string, err error) []string {
	ec.Check(err)
	return value
}

// SyncStrategy provides the string part of the given fallible function result
// while registering the given error.
func (ec *ErrorChecker) SyncStrategy(value config.SyncStrategy, err error) config.SyncStrategy {
	ec.Check(err)
	return value
}
