package runstate

import (
	"fmt"
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
