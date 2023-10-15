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
		if IsCheckoutOpcode(opcode) {
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

// OpcodeTypes provides the names of the types of the opcodes in this program.
func (self *Program) OpcodeTypes() []string {
	result := make([]string, len(self.Opcodes))
	for o, opcode := range self.Opcodes {
		result[o] = reflect.TypeOf(opcode).String()
	}
	return result
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

// WrapOptions represents the options given to Wrap.
type WrapOptions struct {
	RunInGitRoot     bool
	StashOpenChanges bool
	MainBranch       domain.LocalBranchName
	InitialBranch    domain.LocalBranchName
	PreviousBranch   domain.LocalBranchName
}

// Wrap wraps the list with opcodes that
// change to the Git root directory or stash away open changes.
// TODO: only wrap if the list actually contains any opcodes.
func (self *Program) Wrap(options WrapOptions) {
	if !options.PreviousBranch.IsEmpty() {
		self.Add(&opcode.PreserveCheckoutHistory{
			InitialBranch:                     options.InitialBranch,
			InitialPreviouslyCheckedOutBranch: options.PreviousBranch,
			MainBranch:                        options.MainBranch,
		})
	}
	if options.StashOpenChanges {
		self.Prepend(&opcode.StashOpenChanges{})
		self.Add(&opcode.RestoreOpenChanges{})
	}
}
