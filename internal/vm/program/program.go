package program

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Program is a mutable collection of Opcodes.
type Program []shared.Opcode

// Append adds the given opcode to the end of this program.
func (self *Program) Add(opcode ...shared.Opcode) {
	*self = append(*self, opcode...)
}

// AppendProgram adds all elements of the given Program to the end of this Program.
func (self *Program) AddProgram(otherProgram Program) {
	*self = append(*self, otherProgram...)
}

// IsEmpty returns whether or not this Program has any elements.
func (self *Program) IsEmpty() bool {
	return len(*self) == 0
}

// MarshalJSON marshals this program to JSON.
func (self Program) MarshalJSON() ([]byte, error) {
	jsonOpcodes := make([]JSON, len(self))
	for o, opcode := range self {
		jsonOpcodes[o] = JSON{Opcode: opcode}
	}
	return json.Marshal(jsonOpcodes)
}

// OpcodeTypes provides the names of the types of the opcodes in this program.
func (self *Program) OpcodeTypes() []string {
	result := make([]string, len(*self))
	for o, opcode := range *self {
		result[o] = reflect.TypeOf(opcode).String()
	}
	return result
}

// Pop removes and provides the first element of this program.
func (self *Program) Pop() Option[shared.Opcode] {
	if self.IsEmpty() {
		return None[shared.Opcode]()
	}
	result := (*self)[0]
	*self = (*self)[1:]
	return Some(result)
}

// Prepend adds the given opcode to the beginning of this program.
func (self *Program) Prepend(other ...shared.Opcode) {
	if len(other) > 0 {
		result := other
		result = append(result, (*self)...)
		*self = result
	}
}

// PrependProgram adds all elements of the given program to the start of this program.
func (self *Program) PrependProgram(otherProgram Program) {
	result := otherProgram
	result = append(result, (*self)...)
	*self = result
}

func (self *Program) RemoveAllButLast(removeType string) Program {
	allIndexes := slice.FindAll(self.OpcodeTypes(), removeType)
	indexesToRemove := slice.TruncateLast(allIndexes)
	return slice.RemoveAt(*self, indexesToRemove...)
}

// Implementation of the fmt.Stringer interface.
func (self *Program) String() string {
	return self.StringIndented("")
}

func (self *Program) StringIndented(indent string) string {
	if len(*self) == 0 {
		return "(empty program)\n"
	}
	sb := strings.Builder{}
	sb.WriteString("Program:\n")
	for o, opcode := range *self {
		sb.WriteString(fmt.Sprintf("%s%d: %#v\n", indent, o+1, opcode))
	}
	return sb.String()
}

func (self *Program) TouchedBranches() []gitdomain.BranchName {
	var result []gitdomain.BranchName
	for _, opcode := range *self {
		result = append(result, shared.BranchesInOpcode(opcode)...)
	}
	return result
}

// UnmarshalJSON unmarshals the program from JSON.
func (self *Program) UnmarshalJSON(b []byte) error {
	var jsonOpcodes []JSON
	if err := json.Unmarshal(b, &jsonOpcodes); err != nil {
		return err
	}
	if len(jsonOpcodes) > 0 {
		*self = make([]shared.Opcode, len(jsonOpcodes))
		for j, jsonOpcode := range jsonOpcodes {
			(*self)[j] = jsonOpcode.Opcode
		}
	}
	return nil
}
