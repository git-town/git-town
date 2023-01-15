package formatters

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cucumber/messages-go/v16"

	"github.com/cucumber/godog/formatters"
	"github.com/cucumber/godog/internal/utils"
)

const nanoSec = 1000000
const spec = "0.1.0"

func init() {
	formatters.Format("events", fmt.Sprintf("Produces JSON event stream, based on spec: %s.", spec), EventsFormatterFunc)
}

// EventsFormatterFunc implements the FormatterFunc for the events formatter
func EventsFormatterFunc(suite string, out io.Writer) formatters.Formatter {
	return &Events{Base: NewBase(suite, out)}
}

// Events - Events formatter
type Events struct {
	*Base
}

func (f *Events) event(ev interface{}) {
	data, err := json.Marshal(ev)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal stream event: %+v - %v", ev, err))
	}
	fmt.Fprintln(f.out, string(data))
}

// Pickle receives scenario.
func (f *Events) Pickle(pickle *messages.Pickle) {
	f.Base.Pickle(pickle)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.event(&struct {
		Event     string `json:"event"`
		Location  string `json:"location"`
		Timestamp int64  `json:"timestamp"`
	}{
		"TestCaseStarted",
		f.scenarioLocation(pickle),
		utils.TimeNowFunc().UnixNano() / nanoSec,
	})

	if len(pickle.Steps) == 0 {
		// @TODO: is status undefined or passed? when there are no steps
		// for this scenario
		f.event(&struct {
			Event     string `json:"event"`
			Location  string `json:"location"`
			Timestamp int64  `json:"timestamp"`
			Status    string `json:"status"`
		}{
			"TestCaseFinished",
			f.scenarioLocation(pickle),
			utils.TimeNowFunc().UnixNano() / nanoSec,
			"undefined",
		})
	}
}

// TestRunStarted is triggered on test start.
func (f *Events) TestRunStarted() {
	f.Base.TestRunStarted()

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.event(&struct {
		Event     string `json:"event"`
		Version   string `json:"version"`
		Timestamp int64  `json:"timestamp"`
		Suite     string `json:"suite"`
	}{
		"TestRunStarted",
		spec,
		utils.TimeNowFunc().UnixNano() / nanoSec,
		f.suiteName,
	})
}

// Feature receives gherkin document.
func (f *Events) Feature(ft *messages.GherkinDocument, p string, c []byte) {
	f.Base.Feature(ft, p, c)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.event(&struct {
		Event    string `json:"event"`
		Location string `json:"location"`
		Source   string `json:"source"`
	}{
		"TestSource",
		fmt.Sprintf("%s:%d", p, ft.Feature.Location.Line),
		string(c),
	})
}

// Summary pushes summary information to JSON stream.
func (f *Events) Summary() {
	// @TODO: determine status
	status := passed

	f.Storage.MustGetPickleStepResultsByStatus(failed)

	if len(f.Storage.MustGetPickleStepResultsByStatus(failed)) > 0 {
		status = failed
	} else if len(f.Storage.MustGetPickleStepResultsByStatus(passed)) == 0 {
		if len(f.Storage.MustGetPickleStepResultsByStatus(undefined)) > len(f.Storage.MustGetPickleStepResultsByStatus(pending)) {
			status = undefined
		} else {
			status = pending
		}
	}

	snips := f.Snippets()
	if len(snips) > 0 {
		snips = "You can implement step definitions for undefined steps with these snippets:\n" + snips
	}

	f.event(&struct {
		Event     string `json:"event"`
		Status    string `json:"status"`
		Timestamp int64  `json:"timestamp"`
		Snippets  string `json:"snippets"`
		Memory    string `json:"memory"`
	}{
		"TestRunFinished",
		status.String(),
		utils.TimeNowFunc().UnixNano() / nanoSec,
		snips,
		"", // @TODO not sure that could be correctly implemented
	})
}

func (f *Events) step(pickle *messages.Pickle, pickleStep *messages.PickleStep) {
	feature := f.Storage.MustGetFeature(pickle.Uri)
	pickleStepResult := f.Storage.MustGetPickleStepResult(pickleStep.Id)
	step := feature.FindStep(pickleStep.AstNodeIds[0])

	var errMsg string
	if pickleStepResult.Err != nil {
		errMsg = pickleStepResult.Err.Error()
	}
	f.event(&struct {
		Event     string `json:"event"`
		Location  string `json:"location"`
		Timestamp int64  `json:"timestamp"`
		Status    string `json:"status"`
		Summary   string `json:"summary,omitempty"`
	}{
		"TestStepFinished",
		fmt.Sprintf("%s:%d", pickle.Uri, step.Location.Line),
		utils.TimeNowFunc().UnixNano() / nanoSec,
		pickleStepResult.Status.String(),
		errMsg,
	})

	if isLastStep(pickle, pickleStep) {
		var status string

		pickleStepResults := f.Storage.MustGetPickleStepResultsByPickleID(pickle.Id)
		for _, stepResult := range pickleStepResults {
			switch stepResult.Status {
			case passed, failed, undefined, pending:
				status = stepResult.Status.String()
			}
		}

		f.event(&struct {
			Event     string `json:"event"`
			Location  string `json:"location"`
			Timestamp int64  `json:"timestamp"`
			Status    string `json:"status"`
		}{
			"TestCaseFinished",
			f.scenarioLocation(pickle),
			utils.TimeNowFunc().UnixNano() / nanoSec,
			status,
		})
	}
}

// Defined receives step definition.
func (f *Events) Defined(pickle *messages.Pickle, pickleStep *messages.PickleStep, def *formatters.StepDefinition) {
	f.Base.Defined(pickle, pickleStep, def)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	feature := f.Storage.MustGetFeature(pickle.Uri)
	step := feature.FindStep(pickleStep.AstNodeIds[0])

	if def != nil {
		matchedDef := f.Storage.MustGetStepDefintionMatch(pickleStep.AstNodeIds[0])

		m := def.Expr.FindStringSubmatchIndex(pickleStep.Text)[2:]
		var args [][2]int
		for i := 0; i < len(m)/2; i++ {
			pair := m[i : i*2+2]
			var idxs [2]int
			idxs[0] = pair[0]
			idxs[1] = pair[1]
			args = append(args, idxs)
		}

		if len(args) == 0 {
			args = make([][2]int, 0)
		}

		f.event(&struct {
			Event    string   `json:"event"`
			Location string   `json:"location"`
			DefID    string   `json:"definition_id"`
			Args     [][2]int `json:"arguments"`
		}{
			"StepDefinitionFound",
			fmt.Sprintf("%s:%d", pickle.Uri, step.Location.Line),
			DefinitionID(matchedDef),
			args,
		})
	}

	f.event(&struct {
		Event     string `json:"event"`
		Location  string `json:"location"`
		Timestamp int64  `json:"timestamp"`
	}{
		"TestStepStarted",
		fmt.Sprintf("%s:%d", pickle.Uri, step.Location.Line),
		utils.TimeNowFunc().UnixNano() / nanoSec,
	})
}

// Passed captures passed step.
func (f *Events) Passed(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Passed(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(pickle, step)
}

// Skipped captures skipped step.
func (f *Events) Skipped(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Skipped(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(pickle, step)
}

// Undefined captures undefined step.
func (f *Events) Undefined(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Undefined(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(pickle, step)
}

// Failed captures failed step.
func (f *Events) Failed(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition, err error) {
	f.Base.Failed(pickle, step, match, err)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(pickle, step)
}

// Pending captures pending step.
func (f *Events) Pending(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Pending(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(pickle, step)
}

func (f *Events) scenarioLocation(pickle *messages.Pickle) string {
	feature := f.Storage.MustGetFeature(pickle.Uri)
	scenario := feature.FindScenario(pickle.AstNodeIds[0])

	line := scenario.Location.Line
	if len(pickle.AstNodeIds) == 2 {
		_, row := feature.FindExample(pickle.AstNodeIds[1])
		line = row.Location.Line
	}

	return fmt.Sprintf("%s:%d", pickle.Uri, line)
}

func isLastStep(pickle *messages.Pickle, step *messages.PickleStep) bool {
	return pickle.Steps[len(pickle.Steps)-1].Id == step.Id
}
