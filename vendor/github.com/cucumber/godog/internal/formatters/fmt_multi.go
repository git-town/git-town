package formatters

import (
	"io"

	"github.com/cucumber/godog/formatters"
	"github.com/cucumber/godog/internal/storage"
	"github.com/cucumber/messages-go/v16"
)

// MultiFormatter passes test progress to multiple formatters.
type MultiFormatter struct {
	formatters []formatter
	repeater   repeater
}

type formatter struct {
	fmt   formatters.FormatterFunc
	out   io.Writer
	close bool
}

type repeater []formatters.Formatter

type storageFormatter interface {
	SetStorage(s *storage.Storage)
}

// SetStorage passes storage to all added formatters.
func (r repeater) SetStorage(s *storage.Storage) {
	for _, f := range r {
		if ss, ok := f.(storageFormatter); ok {
			ss.SetStorage(s)
		}
	}
}

// TestRunStarted triggers TestRunStarted for all added formatters.
func (r repeater) TestRunStarted() {
	for _, f := range r {
		f.TestRunStarted()
	}
}

// Feature triggers Feature for all added formatters.
func (r repeater) Feature(document *messages.GherkinDocument, s string, bytes []byte) {
	for _, f := range r {
		f.Feature(document, s, bytes)
	}
}

// Pickle triggers Pickle for all added formatters.
func (r repeater) Pickle(pickle *messages.Pickle) {
	for _, f := range r {
		f.Pickle(pickle)
	}
}

// Defined triggers Defined for all added formatters.
func (r repeater) Defined(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	for _, f := range r {
		f.Defined(pickle, step, definition)
	}
}

// Failed triggers Failed for all added formatters.
func (r repeater) Failed(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition, err error) {
	for _, f := range r {
		f.Failed(pickle, step, definition, err)
	}
}

// Passed triggers Passed for all added formatters.
func (r repeater) Passed(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	for _, f := range r {
		f.Passed(pickle, step, definition)
	}
}

// Skipped triggers Skipped for all added formatters.
func (r repeater) Skipped(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	for _, f := range r {
		f.Skipped(pickle, step, definition)
	}
}

// Undefined triggers Undefined for all added formatters.
func (r repeater) Undefined(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	for _, f := range r {
		f.Undefined(pickle, step, definition)
	}
}

// Pending triggers Pending for all added formatters.
func (r repeater) Pending(pickle *messages.Pickle, step *messages.PickleStep, definition *formatters.StepDefinition) {
	for _, f := range r {
		f.Pending(pickle, step, definition)
	}
}

// Summary triggers Summary for all added formatters.
func (r repeater) Summary() {
	for _, f := range r {
		f.Summary()
	}
}

// Add adds formatter with output writer.
func (m *MultiFormatter) Add(name string, out io.Writer) {
	f := formatters.FindFmt(name)
	if f == nil {
		panic("formatter not found: " + name)
	}

	m.formatters = append(m.formatters, formatter{
		fmt: f,
		out: out,
	})
}

// FormatterFunc implements the FormatterFunc for the multi formatter.
func (m *MultiFormatter) FormatterFunc(suite string, out io.Writer) formatters.Formatter {
	for _, f := range m.formatters {
		out := out
		if f.out != nil {
			out = f.out
		}

		m.repeater = append(m.repeater, f.fmt(suite, out))
	}

	return m.repeater
}
