package steps

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/step"
)

// List is a collection of Step instances.
// Only use a list if you need the advanced features of this struct.
// If all you need is an immutable list of steps, using a []step.Step is sufficient.
//
//nolint:musttag // StepList is manually serialized, see the `MarshalJSON` method below
type List struct {
	List []step.Step `exhaustruct:"optional"`
}

// NewStepList provides a StepList instance containing the given step.
func NewStepList(initialStep step.Step) List {
	return List{
		List: []step.Step{initialStep},
	}
}

// Append adds the given step to the end of this StepList.
func (stepList *List) Add(step ...step.Step) {
	stepList.List = append(stepList.List, step...)
}

// AppendList adds all elements of the given StepList to the end of this StepList.
func (stepList *List) AddList(otherList List) {
	stepList.List = append(stepList.List, otherList.List...)
}

// IsEmpty returns whether or not this StepList has any elements.
func (stepList *List) IsEmpty() bool {
	return len(stepList.List) == 0
}

// MarshalJSON marshals the step list to JSON.
func (stepList *List) MarshalJSON() ([]byte, error) {
	jsonSteps := make([]JSON, len(stepList.List))
	for s, step := range stepList.List {
		jsonSteps[s] = JSON{Step: step}
	}
	return json.Marshal(jsonSteps)
}

// Peek provides the first element of this StepList.
func (stepList *List) Peek() step.Step { //nolint:ireturn
	if stepList.IsEmpty() {
		return nil
	}
	return stepList.List[0]
}

// Pop removes and provides the first element of this StepList.
func (stepList *List) Pop() step.Step { //nolint:ireturn
	if stepList.IsEmpty() {
		return nil
	}
	result := stepList.List[0]
	stepList.List = stepList.List[1:]
	return result
}

// Prepend adds the given step to the beginning of this StepList.
func (stepList *List) Prepend(other ...step.Step) {
	if len(other) > 0 {
		stepList.List = append(other, stepList.List...)
	}
}

// PrependList adds all elements of the given StepList to the start of this StepList.
func (stepList *List) PrependList(otherList List) {
	stepList.List = append(otherList.List, stepList.List...)
}

func (stepList *List) RemoveAllButLast(removeType string) {
	typeList := stepList.StepTypes()
	occurrences := slice.FindAll(typeList, removeType)
	occurrencesToRemove := slice.TruncateLast(occurrences)
	for o := len(occurrencesToRemove) - 1; o >= 0; o-- {
		stepList.List = slice.RemoveAt(stepList.List, occurrencesToRemove[o])
	}
}

// RemoveDuplicateCheckoutSteps provides this StepList with checkout steps that immediately follow each other removed.
func (stepList *List) RemoveDuplicateCheckoutSteps() List {
	result := make([]step.Step, 0, len(stepList.List))
	// this one is populated only if the last step is a checkout step
	var lastStep step.Step
	for _, step := range stepList.List {
		if IsCheckoutStep(step) {
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
	return List{List: result}
}

// Implementation of the fmt.Stringer interface.
func (stepList *List) String() string {
	return stepList.StringIndented("")
}

func (stepList *List) StringIndented(indent string) string {
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
func (stepList *List) StepTypes() []string {
	result := make([]string, len(stepList.List))
	for s, step := range stepList.List {
		result[s] = reflect.TypeOf(step).String()
	}
	return result
}

// UnmarshalJSON unmarshals the step list from JSON.
func (stepList *List) UnmarshalJSON(b []byte) error {
	var jsonSteps []JSON
	err := json.Unmarshal(b, &jsonSteps)
	if err != nil {
		return err
	}
	if len(jsonSteps) > 0 {
		stepList.List = make([]step.Step, len(jsonSteps))
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
// TODO: only wrap if the list actually contains any steps.
func (stepList *List) Wrap(options WrapOptions) {
	if !options.PreviousBranch.IsEmpty() {
		stepList.Add(&step.PreserveCheckoutHistory{
			InitialBranch:                     options.InitialBranch,
			InitialPreviouslyCheckedOutBranch: options.PreviousBranch,
			MainBranch:                        options.MainBranch,
		})
	}
	if options.StashOpenChanges {
		stepList.Prepend(&step.StashOpenChanges{})
		stepList.Add(&step.RestoreOpenChanges{})
	}
}
