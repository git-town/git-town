package formatters

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/cucumber/messages-go/v16"

	"github.com/cucumber/godog/formatters"
)

func init() {
	formatters.Format("progress", "Prints a character per step.", ProgressFormatterFunc)
}

// ProgressFormatterFunc implements the FormatterFunc for the progress formatter.
func ProgressFormatterFunc(suite string, out io.Writer) formatters.Formatter {
	return NewProgress(suite, out)
}

// NewProgress creates a new progress formatter.
func NewProgress(suite string, out io.Writer) *Progress {
	steps := 0
	return &Progress{
		Base:        NewBase(suite, out),
		StepsPerRow: 70,
		Steps:       &steps,
	}
}

// Progress is a minimalistic formatter.
type Progress struct {
	*Base
	StepsPerRow int
	Steps       *int
}

// Summary renders summary information.
func (f *Progress) Summary() {
	left := math.Mod(float64(*f.Steps), float64(f.StepsPerRow))
	if left != 0 {
		if *f.Steps > f.StepsPerRow {
			fmt.Fprintf(f.out, s(f.StepsPerRow-int(left))+fmt.Sprintf(" %d\n", *f.Steps))
		} else {
			fmt.Fprintf(f.out, " %d\n", *f.Steps)
		}
	}

	var failedStepsOutput []string

	failedSteps := f.Storage.MustGetPickleStepResultsByStatus(failed)
	sort.Sort(sortPickleStepResultsByPickleStepID(failedSteps))

	for _, sr := range failedSteps {
		if sr.Status == failed {
			pickle := f.Storage.MustGetPickle(sr.PickleID)
			pickleStep := f.Storage.MustGetPickleStep(sr.PickleStepID)
			feature := f.Storage.MustGetFeature(pickle.Uri)

			sc := feature.FindScenario(pickle.AstNodeIds[0])
			scenarioDesc := fmt.Sprintf("%s: %s", sc.Keyword, pickle.Name)
			scenarioLine := fmt.Sprintf("%s:%d", pickle.Uri, sc.Location.Line)

			step := feature.FindStep(pickleStep.AstNodeIds[0])
			stepDesc := strings.TrimSpace(step.Keyword) + " " + pickleStep.Text
			stepLine := fmt.Sprintf("%s:%d", pickle.Uri, step.Location.Line)

			failedStepsOutput = append(
				failedStepsOutput,
				s(2)+red(scenarioDesc)+blackb(" # "+scenarioLine),
				s(4)+red(stepDesc)+blackb(" # "+stepLine),
				s(6)+red("Error: ")+redb(fmt.Sprintf("%+v", sr.Err)),
				"",
			)
		}
	}

	if len(failedStepsOutput) > 0 {
		fmt.Fprintln(f.out, "\n\n--- "+red("Failed steps:")+"\n")
		fmt.Fprint(f.out, strings.Join(failedStepsOutput, "\n"))
	}
	fmt.Fprintln(f.out, "")

	f.Base.Summary()
}

func (f *Progress) step(pickleStepID string) {
	pickleStepResult := f.Storage.MustGetPickleStepResult(pickleStepID)

	switch pickleStepResult.Status {
	case passed:
		fmt.Fprint(f.out, green("."))
	case skipped:
		fmt.Fprint(f.out, cyan("-"))
	case failed:
		fmt.Fprint(f.out, red("F"))
	case undefined:
		fmt.Fprint(f.out, yellow("U"))
	case pending:
		fmt.Fprint(f.out, yellow("P"))
	}

	*f.Steps++

	if math.Mod(float64(*f.Steps), float64(f.StepsPerRow)) == 0 {
		fmt.Fprintf(f.out, " %d\n", *f.Steps)
	}
}

// Passed captures passed step.
func (f *Progress) Passed(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Passed(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(step.Id)
}

// Skipped captures skipped step.
func (f *Progress) Skipped(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Skipped(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(step.Id)
}

// Undefined captures undefined step.
func (f *Progress) Undefined(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Undefined(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(step.Id)
}

// Failed captures failed step.
func (f *Progress) Failed(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition, err error) {
	f.Base.Failed(pickle, step, match, err)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(step.Id)
}

// Pending captures pending step.
func (f *Progress) Pending(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Pending(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.step(step.Id)
}
