package formatters

import (
	"sync"

	"github.com/cucumber/godog/formatters"
	messages "github.com/cucumber/messages/go/v21"
)

// WrapOnFlush wrap a `formatters.Formatter` in a `formatters.FlushFormatter`, which only
// executes when `Flush` is called
func WrapOnFlush(fmt formatters.Formatter) formatters.FlushFormatter {
	return &onFlushFormatter{
		fmt: fmt,
		fns: make([]func(), 0),
		mu:  &sync.Mutex{},
	}
}

type onFlushFormatter struct {
	fmt formatters.Formatter
	fns []func()
	mu  *sync.Mutex
}

func (o *onFlushFormatter) Pickle(pickle *messages.Pickle) {
	o.fns = append(o.fns, func() {
		o.fmt.Pickle(pickle)
	})
}

func (o *onFlushFormatter) Passed(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	o.fns = append(o.fns, func() {
		o.fmt.Passed(pickle, step, definition)
	})
}

// Ambiguous implements formatters.Formatter.
func (o *onFlushFormatter) Ambiguous(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition, err error) {
	o.fns = append(o.fns, func() {
		o.fmt.Ambiguous(pickle, step, definition, err)
	})
}

// Defined implements formatters.Formatter.
func (o *onFlushFormatter) Defined(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	o.fns = append(o.fns, func() {
		o.fmt.Defined(pickle, step, definition)
	})
}

// Failed implements formatters.Formatter.
func (o *onFlushFormatter) Failed(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition, err error) {
	o.fns = append(o.fns, func() {
		o.fmt.Failed(pickle, step, definition, err)
	})
}

// Feature implements formatters.Formatter.
func (o *onFlushFormatter) Feature(pickle *messages.GherkinDocument, p string, c []byte) {
	o.fns = append(o.fns, func() {
		o.fmt.Feature(pickle, p, c)
	})
}

// Pending implements formatters.Formatter.
func (o *onFlushFormatter) Pending(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	o.fns = append(o.fns, func() {
		o.fmt.Pending(pickle, step, definition)
	})
}

// Skipped implements formatters.Formatter.
func (o *onFlushFormatter) Skipped(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	o.fns = append(o.fns, func() {
		o.fmt.Skipped(pickle, step, definition)
	})
}

// Summary implements formatters.Formatter.
func (o *onFlushFormatter) Summary() {
	o.fns = append(o.fns, func() {
		o.fmt.Summary()
	})
}

// TestRunStarted implements formatters.Formatter.
func (o *onFlushFormatter) TestRunStarted() {
	o.fns = append(o.fns, func() {
		o.fmt.TestRunStarted()
	})
}

// Undefined implements formatters.Formatter.
func (o *onFlushFormatter) Undefined(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	o.fns = append(o.fns, func() {
		o.fmt.Undefined(pickle, step, definition)
	})
}

// Flush the logs.
func (o *onFlushFormatter) Flush() {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, fn := range o.fns {
		fn()
	}
}
