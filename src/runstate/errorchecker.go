package runstate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/hosting"
)

// ErrorChecker helps avoid excessive error checking
// while gathering a larger number of values through fallible operations.
type ErrorChecker struct {
	Err error `exhaustruct:"optional"`
}

// Bool provides the bool part of the given fallible function result
// while registering the given error.
func (ec *ErrorChecker) Bool(value bool, err error) bool {
	ec.Check(err)
	return value
}

// Fail registers an error with the given format string.
func (ec *ErrorChecker) Fail(format string, a ...any) error {
	ec.Check(fmt.Errorf(format, a...))
	return ec.Err
}

// Check registers the given error and indicates whether this ErrorChecker contains an error now.
func (ec *ErrorChecker) Check(err error) bool {
	if err != nil && ec.Err == nil {
		ec.Err = err
	}
	return ec.Err != nil
}

// String provides the string part of the given fallible function result
// while registering the given error.
func (ec *ErrorChecker) Proposal(value *hosting.Proposal, err error) *hosting.Proposal {
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
