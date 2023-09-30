package runstate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/src/steps"
)

// StepList is a fifo containing Step instances.
//
//nolint:musttag // StepList is manually serialized, see the `MarshalJSON` method below
type StepList struct {
	List []steps.Step `exhaustruct:"optional"`
}

// NewStepList provides a StepList instance containing the given step.
func NewStepList(step steps.Step) StepList {
	return StepList{
		List: []steps.Step{step},
	}
}

// Append adds the given step to the end of this StepList.
func (stepList *StepList) Append(step ...steps.Step) {
	stepList.List = append(stepList.List, step...)
}

// AppendList adds all elements of the given StepList to the end of this StepList.
func (stepList *StepList) AppendList(otherList StepList) {
	stepList.List = append(stepList.List, otherList.List...)
}

// IsEmpty returns whether or not this StepList has any elements.
func (stepList *StepList) IsEmpty() bool {
	return len(stepList.List) == 0
}

// MarshalJSON marshals the step list to JSON.
func (stepList *StepList) MarshalJSON() ([]byte, error) {
	jsonSteps := make([]JSONStep, len(stepList.List))
	for s, step := range stepList.List {
		jsonSteps[s] = JSONStep{Step: step}
	}
	return json.Marshal(jsonSteps)
}

// Peek provides the first element of this StepList.
func (stepList *StepList) Peek() steps.Step {
	if stepList.IsEmpty() {
		return nil
	}
	return stepList.List[0]
}

// Pop removes and provides the first element of this StepList.
func (stepList *StepList) Pop() steps.Step {
	if stepList.IsEmpty() {
		return nil
	}
	result := stepList.List[0]
	stepList.List = stepList.List[1:]
	return result
}

// Prepend adds the given step to the beginning of this StepList.
func (stepList *StepList) Prepend(other ...steps.Step) {
	if len(other) > 0 {
		stepList.List = append(other, stepList.List...)
	}
}

// PrependList adds all elements of the given StepList to the start of this StepList.
func (stepList *StepList) PrependList(otherList StepList) {
	stepList.List = append(otherList.List, stepList.List...)
}

func (stepList *StepList) RemoveAllButLast(removeType string) {
	typeList := stepList.StepTypes()
	occurrences := slice.FindAll(typeList, removeType)
	occurrencesToRemove := slice.TruncateLast(occurrences)
	for o := len(occurrencesToRemove) - 1; o >= 0; o-- {
		stepList.List = slice.RemoveAt(stepList.List, occurrencesToRemove[o])
	}
}

// RemoveDuplicateCheckoutSteps provides this StepList with checkout steps that immediately follow each other removed.
func (stepList *StepList) RemoveDuplicateCheckoutSteps() StepList {
	result := make([]steps.Step, 0, len(stepList.List))
	// this one is populated only if the last step is a checkout step
	var lastStep steps.Step
	for _, step := range stepList.List {
		if isCheckoutStep(step) {
			lastStep = step
			continue
		}
		if lastStep != nil {
			result = append(result, lastStep)
		}
		lastStep = nil
		result = append(result, step)
	}
	if lastStep != nil {
		result = append(result, lastStep)
	}
	return StepList{List: result}
}

// Implementation of the fmt.Stringer interface.
func (stepList *StepList) String() string {
	return stepList.StringIndented("")
}

func (stepList *StepList) StringIndented(indent string) string {
	sb := strings.Builder{}
	if stepList.IsEmpty() {
		sb.WriteString("(empty StepList)\n")
	} else {
		sb.WriteString("StepList:\n")
		for s, step := range stepList.List {
			sb.WriteString(fmt.Sprintf("%s%d: %#v\n", indent, s+1, step))
		}
	}
	return sb.String()
}

// StepTypes provides the names of the types of the steps in this list.
func (stepList *StepList) StepTypes() []string {
	result := make([]string, len(stepList.List))
	for s, step := range stepList.List {
		result[s] = reflect.TypeOf(step).String()
	}
	return result
}

// UnmarshalJSON unmarshals the step list from JSON.
func (stepList *StepList) UnmarshalJSON(b []byte) error {
	var jsonSteps []JSONStep
	err := json.Unmarshal(b, &jsonSteps)
	if err != nil {
		return err
	}
	if len(jsonSteps) > 0 {
		stepList.List = make([]steps.Step, len(jsonSteps))
		for j, jsonStep := range jsonSteps {
			stepList.List[j] = jsonStep.Step
		}
	}
	return nil
}

// WrapOptions represents the options given to Wrap.
type WrapOptions struct {
	RunInGitRoot     bool
	StashOpenChanges bool
	MainBranch       domain.LocalBranchName
	InitialBranch    domain.LocalBranchName
	PreviousBranch   domain.LocalBranchName
}

// Wrap wraps the list with steps that
// change to the Git root directory or stash away open changes.
func (stepList *StepList) Wrap(options WrapOptions) error {
	if !options.PreviousBranch.IsEmpty() {
		stepList.Append(&steps.PreserveCheckoutHistoryStep{
			InitialBranch:                     options.InitialBranch,
			InitialPreviouslyCheckedOutBranch: options.PreviousBranch,
			MainBranch:                        options.MainBranch,
		})
	}
	if options.StashOpenChanges {
		stepList.Prepend(&steps.StashOpenChangesStep{})
		stepList.Append(&steps.RestoreOpenChangesStep{})
	}
	return nil
}
