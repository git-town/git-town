package program

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// Program is a collection of Step instances.
// Only use a list if you need the advanced features of this struct.
// If all you need is an immutable list of steps, using a []step.Step is sufficient.
//
//nolint:musttag // program is manually serialized, see the `MarshalJSON` method below
type Program struct {
	Opcodes []shared.Opcode `exhaustruct:"optional"`
}

// NewProgram provides a program instance containing the given step.
func NewProgram(initialStep shared.Opcode) Program {
	return Program{
		Opcodes: []shared.Opcode{initialStep},
	}
}

// Append adds the given step to the end of this program.
func (p *Program) Add(step ...shared.Opcode) {
	p.Opcodes = append(p.Opcodes, step...)
}

// AppendProgram adds all elements of the given Program to the end of this Program.
func (p *Program) AddProgram(otherProgram Program) {
	p.Opcodes = append(p.Opcodes, otherProgram.Opcodes...)
}

// IsEmpty returns whether or not this Program has any elements.
func (p *Program) IsEmpty() bool {
	return len(p.Opcodes) == 0
}

// MarshalJSON marshals the step list to JSON.
func (p *Program) MarshalJSON() ([]byte, error) {
	jsonSteps := make([]JSON, len(p.Opcodes))
	for s, step := range p.Opcodes {
		jsonSteps[s] = JSON{Opcode: step}
	}
	return json.Marshal(jsonSteps)
}

// Peek provides the first element of this program.
func (p *Program) Peek() shared.Opcode { //nolint:ireturn
	if p.IsEmpty() {
		return nil
	}
	return p.Opcodes[0]
}

// Pop removes and provides the first element of this program.
func (p *Program) Pop() shared.Opcode { //nolint:ireturn
	if p.IsEmpty() {
		return nil
	}
	result := p.Opcodes[0]
	p.Opcodes = p.Opcodes[1:]
	return result
}

// Prepend adds the given step to the beginning of this program.
func (p *Program) Prepend(other ...shared.Opcode) {
	if len(other) > 0 {
		p.Opcodes = append(other, p.Opcodes...)
	}
}

// PrependProgram adds all elements of the given program to the start of this program.
func (p *Program) PrependProgram(otherProgram Program) {
	p.Opcodes = append(otherProgram.Opcodes, p.Opcodes...)
}

func (p *Program) RemoveAllButLast(removeType string) {
	stepTypes := p.StepTypes()
	occurrences := slice.FindAll(stepTypes, removeType)
	occurrencesToRemove := slice.TruncateLast(occurrences)
	for o := len(occurrencesToRemove) - 1; o >= 0; o-- {
		p.Opcodes = slice.RemoveAt(p.Opcodes, occurrencesToRemove[o])
	}
}

// RemoveDuplicateCheckoutSteps provides this program with checkout steps that immediately follow each other removed.
func (p *Program) RemoveDuplicateCheckoutSteps() Program {
	result := make([]shared.Opcode, 0, len(p.Opcodes))
	// this one is populated only if the last step is a checkout step
	var lastStep shared.Opcode
	for _, step := range p.Opcodes {
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
	return Program{Opcodes: result}
}

// Implementation of the fmt.Stringer interface.
func (p *Program) String() string {
	return p.StringIndented("")
}

func (p *Program) StringIndented(indent string) string {
	sb := strings.Builder{}
	if p.IsEmpty() {
		sb.WriteString("(empty program)\n")
	} else {
		sb.WriteString("Program:\n")
		for s, step := range p.Opcodes {
			sb.WriteString(fmt.Sprintf("%s%d: %#v\n", indent, s+1, step))
		}
	}
	return sb.String()
}

// StepTypes provides the names of the types of the steps in this list.
func (p *Program) StepTypes() []string {
	result := make([]string, len(p.Opcodes))
	for s, step := range p.Opcodes {
		result[s] = reflect.TypeOf(step).String()
	}
	return result
}

// UnmarshalJSON unmarshals the step list from JSON.
func (p *Program) UnmarshalJSON(b []byte) error {
	var jsonSteps []JSON
	err := json.Unmarshal(b, &jsonSteps)
	if err != nil {
		return err
	}
	if len(jsonSteps) > 0 {
		p.Opcodes = make([]shared.Opcode, len(jsonSteps))
		for j, jsonStep := range jsonSteps {
			p.Opcodes[j] = jsonStep.Opcode
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
func (p *Program) Wrap(options WrapOptions) {
	if !options.PreviousBranch.IsEmpty() {
		p.Add(&opcode.PreserveCheckoutHistory{
			InitialBranch:                     options.InitialBranch,
			InitialPreviouslyCheckedOutBranch: options.PreviousBranch,
			MainBranch:                        options.MainBranch,
		})
	}
	if options.StashOpenChanges {
		p.Prepend(&opcode.StashOpenChanges{})
		p.Add(&opcode.RestoreOpenChanges{})
	}
}
