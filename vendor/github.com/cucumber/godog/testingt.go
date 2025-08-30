package godog

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// T returns a TestingT compatible interface from the current test context. It will return nil if
// called outside the context of a test. This can be used with (for example) testify's assert and
// require packages.
func T(ctx context.Context) TestingT {
	return getTestingT(ctx)
}

// TestingT is a subset of the public methods implemented by go's testing.T. It allows assertion
// libraries to be used with godog, provided they depend only on this subset of methods.
type TestingT interface {
	// Name returns the name of the current pickle under test
	Name() string
	// Log will log to the current testing.T log if set, otherwise it will log to stdout
	Log(args ...interface{})
	// Logf will log a formatted string to the current testing.T log if set, otherwise it will log
	// to stdout
	Logf(format string, args ...interface{})
	// Error fails the current test and logs the provided arguments. Equivalent to calling Log then
	// Fail.
	Error(args ...interface{})
	// Errorf fails the current test and logs the formatted message. Equivalent to calling Logf then
	// Fail.
	Errorf(format string, args ...interface{})
	// Fail marks the current test as failed, but does not halt execution of the step.
	Fail()
	// FailNow marks the current test as failed and halts execution of the step.
	FailNow()
	// Fatal logs the provided arguments, marks the test as failed and halts execution of the step.
	Fatal(args ...interface{})
	// Fatal logs the formatted message, marks the test as failed and halts execution of the step.
	Fatalf(format string, args ...interface{})
	// Skip logs the provided arguments and marks the test as skipped but does not halt execution
	// of the step.
	Skip(args ...interface{})
	// Skipf logs the formatted message and marks the test as skipped but does not halt execution
	// of the step.
	Skipf(format string, args ...interface{})
	// SkipNow marks the current test as skipped and halts execution of the step.
	SkipNow()
	// Skipped returns true if the test has been marked as skipped.
	Skipped() bool
}

// Logf will log test output. If called in the context of a test and testing.T has been registered,
// this will log using the step's testing.T, else it will simply log to stdout.
func Logf(ctx context.Context, format string, args ...interface{}) {
	if t := getTestingT(ctx); t != nil {
		t.Logf(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

// Log will log test output. If called in the context of a test and testing.T has been registered,
// this will log using the step's testing.T, else it will simply log to stdout.
func Log(ctx context.Context, args ...interface{}) {
	if t := getTestingT(ctx); t != nil {
		t.Log(args...)
		return
	}
	fmt.Println(args...)
}

// LoggedMessages returns an array of any logged messages that have been recorded during the test
// through calls to godog.Log / godog.Logf or via operations against godog.T(ctx)
func LoggedMessages(ctx context.Context) []string {
	if t := getTestingT(ctx); t != nil {
		return t.logMessages
	}
	return nil
}

// errStopNow should be returned inside a panic within the test to immediately halt execution of that
// test
var errStopNow = fmt.Errorf("FailNow or SkipNow called")

type testingT struct {
	name         string
	t            *testing.T
	failed       bool
	skipped      bool
	failMessages []string
	logMessages  []string
}

// check interface against our testingT and the upstream testing.B/F/T:
var (
	_ TestingT = &testingT{}
	_ TestingT = (*testing.T)(nil)
)

func (dt *testingT) Name() string {
	if dt.t != nil {
		return dt.t.Name()
	}
	return dt.name
}

func (dt *testingT) Log(args ...interface{}) {
	dt.logMessages = append(dt.logMessages, fmt.Sprint(args...))
	if dt.t != nil {
		dt.t.Log(args...)
		return
	}
	fmt.Println(args...)
}

func (dt *testingT) Logf(format string, args ...interface{}) {
	dt.logMessages = append(dt.logMessages, fmt.Sprintf(format, args...))
	if dt.t != nil {
		dt.t.Logf(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

func (dt *testingT) Error(args ...interface{}) {
	dt.Log(args...)
	dt.failMessages = append(dt.failMessages, fmt.Sprintln(args...))
	dt.Fail()
}

func (dt *testingT) Errorf(format string, args ...interface{}) {
	dt.Logf(format, args...)
	dt.failMessages = append(dt.failMessages, fmt.Sprintf(format, args...))
	dt.Fail()
}

func (dt *testingT) Fail() {
	dt.failed = true
}

func (dt *testingT) FailNow() {
	dt.Fail()
	panic(errStopNow)
}

func (dt *testingT) Fatal(args ...interface{}) {
	dt.Log(args...)
	dt.FailNow()
}

func (dt *testingT) Fatalf(format string, args ...interface{}) {
	dt.Logf(format, args...)
	dt.FailNow()
}

func (dt *testingT) Skip(args ...interface{}) {
	dt.Log(args...)
	dt.skipped = true
}

func (dt *testingT) Skipf(format string, args ...interface{}) {
	dt.Logf(format, args...)
	dt.skipped = true
}

func (dt *testingT) SkipNow() {
	dt.skipped = true
	panic(errStopNow)
}

func (dt *testingT) Skipped() bool {
	return dt.skipped
}

// isFailed will return an error representing the calls to Fail made during this test
func (dt *testingT) isFailed() error {
	if dt.skipped {
		return ErrSkip
	}
	if !dt.failed {
		return nil
	}
	switch len(dt.failMessages) {
	case 0:
		return fmt.Errorf("fail called on TestingT")
	case 1:
		return fmt.Errorf(dt.failMessages[0])
	default:
		return fmt.Errorf("checks failed:\n* %s", strings.Join(dt.failMessages, "\n* "))
	}
}

type testingTCtxVal struct{}

func setContextTestingT(ctx context.Context, dt *testingT) context.Context {
	return context.WithValue(ctx, testingTCtxVal{}, dt)
}

func getTestingT(ctx context.Context) *testingT {
	dt, ok := ctx.Value(testingTCtxVal{}).(*testingT)
	if !ok {
		return nil
	}
	return dt
}
