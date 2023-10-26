// Package opcode defines the individual operations that the Git Town VM can execute.
// All opcodes implement the shared.Opcode interface.
package opcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// undeclaredOpcodeMethods can be added to structs in this package to satisfy the shared.Opcode interface even if they don't declare all required methods.
type undeclaredOpcodeMethods struct{}

func (self *undeclaredOpcodeMethods) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (self *undeclaredOpcodeMethods) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) Run(_ shared.RunArgs) error {
	return nil
}

func (self *undeclaredOpcodeMethods) ShouldAutomaticallyAbortOnError() bool {
	return false
}

func Lookup(opcodeType string) shared.Opcode { //nolint:ireturn
	for _, opcode := range Types() {
		if gohacks.TypeName(opcode) == opcodeType {
			return opcode
		}
	}
	return nil
}

// Program is a mutable collection of Opcodes.
//
//nolint:musttag // program is manually serialized, see the `MarshalJSON` method below
type Program struct {
	Opcodes []shared.Opcode `exhaustruct:"optional"`
}

// Append adds the given opcode to the end of this program.
func (program *Program) Add(opcode ...shared.Opcode) {
	program.Opcodes = append(program.Opcodes, opcode...)
}

// AppendProgram adds all elements of the given Program to the end of this Program.
func (program *Program) AddProgram(otherProgram Program) {
	program.Opcodes = append(program.Opcodes, otherProgram.Opcodes...)
}

// IsEmpty returns whether or not this Program has any elements.
func (program *Program) IsEmpty() bool {
	return len(program.Opcodes) == 0
}

// MarshalJSON marshals this program to JSON.
func (program *Program) MarshalJSON() ([]byte, error) {
	jsonOpcodes := make([]JSON, len(program.Opcodes))
	for o, opcode := range program.Opcodes {
		jsonOpcodes[o] = JSON{Opcode: opcode}
	}
	return json.Marshal(jsonOpcodes)
}

// OpcodeTypes provides the names of the types of the opcodes in this program.
func (program *Program) OpcodeTypes() []string {
	result := make([]string, len(program.Opcodes))
	for o, opcode := range program.Opcodes {
		result[o] = reflect.TypeOf(opcode).String()
	}
	return result
}

// Peek provides the first element of this program.
func (program *Program) Peek() shared.Opcode { //nolint:ireturn
	if program.IsEmpty() {
		return nil
	}
	return program.Opcodes[0]
}

// Pop removes and provides the first element of this program.
func (program *Program) Pop() shared.Opcode { //nolint:ireturn
	if program.IsEmpty() {
		return nil
	}
	result := program.Opcodes[0]
	program.Opcodes = program.Opcodes[1:]
	return result
}

// Prepend adds the given opcode to the beginning of this program.
func (program *Program) Prepend(other ...shared.Opcode) {
	if len(other) > 0 {
		program.Opcodes = append(other, program.Opcodes...)
	}
}

// PrependProgram adds all elements of the given program to the start of this program.
func (program *Program) PrependProgram(otherProgram Program) {
	program.Opcodes = append(otherProgram.Opcodes, program.Opcodes...)
}

func (program *Program) RemoveAllButLast(removeType string) {
	opcodeTypes := program.OpcodeTypes()
	occurrences := slice.FindAll(opcodeTypes, removeType)
	occurrencesToRemove := slice.TruncateLast(occurrences)
	for o := len(occurrencesToRemove) - 1; o >= 0; o-- {
		program.Opcodes = slice.RemoveAt(program.Opcodes, occurrencesToRemove[o])
	}
}

// RemoveDuplicateCheckout provides this program with checkout opcodes that immediately follow each other removed.
func (program *Program) RemoveDuplicateCheckout() Program {
	result := make([]shared.Opcode, 0, len(program.Opcodes))
	// this one is populated only if the last opcode is a checkout
	var lastOpcode shared.Opcode
	for _, opcode := range program.Opcodes {
		if shared.IsCheckoutOpcode(opcode) {
			lastOpcode = opcode
			continue
		}
		if lastOpcode != nil {
			result = append(result, lastOpcode)
		}
		lastOpcode = nil
		result = append(result, opcode)
	}
	if lastOpcode != nil {
		result = append(result, lastOpcode)
	}
	return Program{Opcodes: result}
}

// Implementation of the fmt.Stringer interface.
func (program *Program) String() string {
	return program.StringIndented("")
}

func (program *Program) StringIndented(indent string) string {
	sb := strings.Builder{}
	if program.IsEmpty() {
		sb.WriteString("(empty program)\n")
	} else {
		sb.WriteString("Program:\n")
		for o, opcode := range program.Opcodes {
			sb.WriteString(fmt.Sprintf("%s%d: %#v\n", indent, o+1, opcode))
		}
	}
	return sb.String()
}

// UnmarshalJSON unmarshals the program from JSON.
func (program *Program) UnmarshalJSON(b []byte) error {
	var jsonOpcodes []JSON
	err := json.Unmarshal(b, &jsonOpcodes)
	if err != nil {
		return err
	}
	if len(jsonOpcodes) > 0 {
		program.Opcodes = make([]shared.Opcode, len(jsonOpcodes))
		for j, jsonOpcode := range jsonOpcodes {
			program.Opcodes[j] = jsonOpcode.Opcode
		}
	}
	return nil
}

// JSON is used to store an opcode in JSON.
type JSON struct { //nolint:musttag // JSON uses a custom serialization algorithm
	Opcode shared.Opcode
}

// MarshalJSON marshals the opcode to JSON.
func (j *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": j.Opcode,
		"type": gohacks.TypeName(j.Opcode),
	})
}

// UnmarshalJSON unmarshals the opcode from JSON.
func (j *JSON) UnmarshalJSON(b []byte) error {
	var mapping map[string]json.RawMessage
	err := json.Unmarshal(b, &mapping)
	if err != nil {
		return err
	}
	var opcodeType string
	err = json.Unmarshal(mapping["type"], &opcodeType)
	if err != nil {
		return err
	}
	j.Opcode = Lookup(opcodeType)
	if j.Opcode == nil {
		return fmt.Errorf(messages.OpcodeUnknown, opcodeType)
	}
	return json.Unmarshal(mapping["data"], &j.Opcode)
}

// Types provides all existing opcodes.
// This is used to iterate all opcode types.
func Types() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
		&AbortRebase{},
		&AddToPerennialBranches{},
		&ChangeParent{},
		&Checkout{},
		&CheckoutIfExists{},
		&CheckoutParent{},
		&ChangeParent{},
		&CommitOpenChanges{},
		&ConnectorMergeProposal{},
		&ContinueMerge{},
		&ContinueRebase{},
		&CreateBranch{},
		&CreateBranchExistingParent{},
		&CreateProposal{},
		&CreateRemoteBranch{},
		&CreateTrackingBranch{},
		&DeleteLocalBranch{},
		&DeleteParentBranch{},
		&DeleteRemoteBranch{},
		&DeleteTrackingBranch{},
		&DiscardOpenChanges{},
		&EnsureHasShippableChanges{},
		&FetchUpstream{},
		&ForcePushCurrentBranch{},
		&IfBranchHasUnmergedChanges{},
		&Merge{},
		&MergeParent{},
		&PreserveCheckoutHistory{},
		&PullCurrentBranch{},
		&PushCurrentBranch{},
		&PushTags{},
		&RebaseBranch{},
		&RebaseParent{},
		&RemoveBranchFromLineage{},
		&RemoveFromPerennialBranches{},
		&RemoveGlobalConfig{},
		&RemoveLocalConfig{},
		&ResetCurrentBranchToSHA{},
		&ResetRemoteBranchToSHA{},
		&RestoreOpenChanges{},
		&RevertCommit{},
		&SetExistingParent{},
		&SetGlobalConfig{},
		&SetLocalConfig{},
		&SetParent{},
		&SkipCurrentBranch{},
		&StashOpenChanges{},
		&SquashMerge{},
		&UndoLastCommit{},
		&UpdateProposalTarget{},
	}
}
