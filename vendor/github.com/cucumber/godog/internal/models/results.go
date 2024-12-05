package models

import (
	"time"

	"github.com/cucumber/godog/colors"
	"github.com/cucumber/godog/internal/utils"
)

// TestRunStarted ...
type TestRunStarted struct {
	StartedAt time.Time
}

// PickleResult ...
type PickleResult struct {
	PickleID  string
	StartedAt time.Time
}

// PickleAttachment ...
type PickleAttachment struct {
	Name     string
	MimeType string
	Data     []byte
}

// PickleStepResult ...
type PickleStepResult struct {
	Status     StepResultStatus
	FinishedAt time.Time
	Err        error

	PickleID     string
	PickleStepID string

	Def *StepDefinition

	Attachments []PickleAttachment
}

// NewStepResult ...
func NewStepResult(
	status StepResultStatus,
	pickleID, pickleStepID string,
	match *StepDefinition,
	attachments []PickleAttachment,
	err error,
) PickleStepResult {
	return PickleStepResult{
		Status:       status,
		FinishedAt:   utils.TimeNowFunc(),
		Err:          err,
		PickleID:     pickleID,
		PickleStepID: pickleStepID,
		Def:          match,
		Attachments:  attachments,
	}
}

// StepResultStatus ...
type StepResultStatus int

const (
	// Passed ...
	Passed StepResultStatus = iota
	// Failed ...
	Failed
	// Skipped ...
	Skipped
	// Undefined ...
	Undefined
	// Pending ...
	Pending
	// Ambiguous ...
	Ambiguous
)

// Color ...
func (st StepResultStatus) Color() colors.ColorFunc {
	switch st {
	case Passed:
		return colors.Green
	case Failed:
		return colors.Red
	case Skipped:
		return colors.Cyan
	default:
		return colors.Yellow
	}
}

// String ...
func (st StepResultStatus) String() string {
	switch st {
	case Passed:
		return "passed"
	case Failed:
		return "failed"
	case Skipped:
		return "skipped"
	case Undefined:
		return "undefined"
	case Pending:
		return "pending"
	case Ambiguous:
		return "ambiguous"
	default:
		return "unknown"
	}
}
