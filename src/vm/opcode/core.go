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
	switch opcodeType {
	case "AbortMerge":
		return &AbortMerge{}
	case "AbortRebase":
		return &AbortRebase{}
	case "AddToPerennialBranches":
		return &AddToPerennialBranches{}
	case "ChangeParent":
		return &ChangeParent{}
	case "Checkout":
		return &Checkout{}
	case "CheckoutIfExists":
		return &CheckoutIfExists{}
	case "CheckoutParent":
		return &CheckoutParent{}
	case "CommitOpenChanges":
		return &CommitOpenChanges{}
	case "ConnectorMergeProposal":
		return &ConnectorMergeProposal{}
	case "ContinueMerge":
		return &ContinueMerge{}
	case "ContinueRebase":
		return &ContinueRebase{}
	case "CreateBranch":
		return &CreateBranch{}
	case "CreateBranchExistingParent":
		return &CreateBranchExistingParent{}
	case "CreateProposal":
		return &CreateProposal{}
	case "CreateRemoteBranch":
		return &CreateRemoteBranch{}
	case "CreateTrackingBranch":
		return &CreateTrackingBranch{}
	case "DeleteLocalBranch":
		return &DeleteLocalBranch{}
	case "DeleteParentBranch":
		return &DeleteParentBranch{}
	case "DeleteRemoteBranch":
		return &DeleteRemoteBranch{}
	case "DeleteTrackingBranch":
		return &DeleteTrackingBranch{}
	case "DiscardOpenChanges":
		return &DiscardOpenChanges{}
	case "EnsureHasShippableChanges":
		return &EnsureHasShippableChanges{}
	case "FetchUpstream":
		return &FetchUpstream{}
	case "ForcePushCurrentBranch":
		return &ForcePushCurrentBranch{}
	case "IfElse":
		return &IfElse{}
	case "Merge":
		return &Merge{}
	case "MergeParent":
		return &MergeParent{}
	case "PreserveCheckoutHistory":
		return &PreserveCheckoutHistory{}
	case "PullCurrentBranch":
		return &PullCurrentBranch{}
	case "PushCurrentBranch":
		return &PushCurrentBranch{}
	case "PushTags":
		return &PushTags{}
	case "RebaseBranch":
		return &RebaseBranch{}
	case "RebaseParent":
		return &RebaseParent{}
	case "RemoveFromPerennialBranches":
		return &RemoveFromPerennialBranches{}
	case "RemoveGlobalConfig":
		return &RemoveGlobalConfig{}
	case "RemoveLocalConfig":
		return &RemoveLocalConfig{}
	case "ResetCurrentBranchToSHA":
		return &ResetCurrentBranchToSHA{}
	case "ResetRemoteBranchToSHA":
		return &ResetRemoteBranchToSHA{}
	case "RestoreOpenChanges":
		return &RestoreOpenChanges{}
	case "RevertCommit":
		return &RevertCommit{}
	case "SetExistingParent":
		return &SetExistingParent{}
	case "SetGlobalConfig":
		return &SetGlobalConfig{}
	case "SetLocalConfig":
		return &SetLocalConfig{}
	case "SetParent":
		return &SetParent{}
	case "SetParentIfBranchExists":
		return &SetParentIfBranchExists{}
	case "SquashMerge":
		return &SquashMerge{}
	case "SkipCurrentBranch":
		return &SkipCurrentBranch{}
	case "StashOpenChanges":
		return &StashOpenChanges{}
	case "UndoLastCommit":
		return &UndoLastCommit{}
	case "UpdateProposalTarget":
		return &UpdateProposalTarget{}
	}
	return nil
}

// Program is a mutable collection of Opcodes.
// Only use a program if you need the mutability features of this struct.
// If all you need is an immutable list of opcodes, a []shared.Opcode is sufficient.
//
//nolint:musttag // program is manually serialized, see the `MarshalJSON` method below
type Program struct {
	Opcodes []shared.Opcode `exhaustruct:"optional"`
}

// Append adds the given opcode to the end of this program.
func (self *Program) Add(opcode ...shared.Opcode) {
	self.Opcodes = append(self.Opcodes, opcode...)
}

// AppendProgram adds all elements of the given Program to the end of this Program.
func (self *Program) AddProgram(otherProgram Program) {
	self.Opcodes = append(self.Opcodes, otherProgram.Opcodes...)
}

// IsEmpty returns whether or not this Program has any elements.
func (self *Program) IsEmpty() bool {
	return len(self.Opcodes) == 0
}

// MarshalJSON marshals this program to JSON.
func (self *Program) MarshalJSON() ([]byte, error) {
	jsonOpcodes := make([]JSON, len(self.Opcodes))
	for o, opcode := range self.Opcodes {
		jsonOpcodes[o] = JSON{Opcode: opcode}
	}
	return json.Marshal(jsonOpcodes)
}

// OpcodeTypes provides the names of the types of the opcodes in this program.
func (self *Program) OpcodeTypes() []string {
	result := make([]string, len(self.Opcodes))
	for o, opcode := range self.Opcodes {
		result[o] = reflect.TypeOf(opcode).String()
	}
	return result
}

// Peek provides the first element of this program.
func (self *Program) Peek() shared.Opcode { //nolint:ireturn
	if self.IsEmpty() {
		return nil
	}
	return self.Opcodes[0]
}

// Pop removes and provides the first element of this program.
func (self *Program) Pop() shared.Opcode { //nolint:ireturn
	if self.IsEmpty() {
		return nil
	}
	result := self.Opcodes[0]
	self.Opcodes = self.Opcodes[1:]
	return result
}

// Prepend adds the given opcode to the beginning of this program.
func (self *Program) Prepend(other ...shared.Opcode) {
	if len(other) > 0 {
		self.Opcodes = append(other, self.Opcodes...)
	}
}

// PrependProgram adds all elements of the given program to the start of this program.
func (self *Program) PrependProgram(otherProgram Program) {
	self.Opcodes = append(otherProgram.Opcodes, self.Opcodes...)
}

func (self *Program) RemoveAllButLast(removeType string) {
	opcodeTypes := self.OpcodeTypes()
	occurrences := slice.FindAll(opcodeTypes, removeType)
	occurrencesToRemove := slice.TruncateLast(occurrences)
	for o := len(occurrencesToRemove) - 1; o >= 0; o-- {
		self.Opcodes = slice.RemoveAt(self.Opcodes, occurrencesToRemove[o])
	}
}

// RemoveDuplicateCheckout provides this program with checkout opcodes that immediately follow each other removed.
func (self *Program) RemoveDuplicateCheckout() Program {
	result := make([]shared.Opcode, 0, len(self.Opcodes))
	// this one is populated only if the last opcode is a checkout
	var lastOpcode shared.Opcode
	for _, opcode := range self.Opcodes {
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
func (self *Program) String() string {
	return self.StringIndented("")
}

func (self *Program) StringIndented(indent string) string {
	sb := strings.Builder{}
	if self.IsEmpty() {
		sb.WriteString("(empty program)\n")
	} else {
		sb.WriteString("Program:\n")
		for o, opcode := range self.Opcodes {
			sb.WriteString(fmt.Sprintf("%s%d: %#v\n", indent, o+1, opcode))
		}
	}
	return sb.String()
}

// UnmarshalJSON unmarshals the program from JSON.
func (self *Program) UnmarshalJSON(b []byte) error {
	var jsonOpcodes []JSON
	err := json.Unmarshal(b, &jsonOpcodes)
	if err != nil {
		return err
	}
	if len(jsonOpcodes) > 0 {
		self.Opcodes = make([]shared.Opcode, len(jsonOpcodes))
		for j, jsonOpcode := range jsonOpcodes {
			self.Opcodes[j] = jsonOpcode.Opcode
		}
	}
	return nil
}

// JSON is used to store an opcode in JSON.
type JSON struct { //nolint:musttag // JSON uses a custom serialization algorithm
	Opcode shared.Opcode
}

// MarshalJSON marshals the opcode to JSON.
func (self *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": self.Opcode,
		"type": gohacks.TypeName(self.Opcode),
	})
}

// UnmarshalJSON unmarshals the opcode from JSON.
func (self *JSON) UnmarshalJSON(b []byte) error {
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
	self.Opcode = Lookup(opcodeType)
	if self.Opcode == nil {
		return fmt.Errorf(messages.OpcodeUnknown, opcodeType)
	}
	return json.Unmarshal(mapping["data"], &self.Opcode)
}
