package program

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/git-town/git-town/v12/src/gohacks/slice"
	"github.com/git-town/git-town/v12/src/vm/shared"
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
func (self Program) IsEmpty() bool {
	return len(self) == 0
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
func (self Program) OpcodeTypes() []string {
	result := make([]string, len(self))
	for o, opcode := range self {
		result[o] = reflect.TypeOf(opcode).String()
	}
	return result
}

// Peek provides the first element of this program.
func (self Program) Peek() shared.Opcode { //nolint:ireturn
	if self.IsEmpty() {
		return nil
	}
	return self[0]
}

// Pop removes and provides the first element of this program.
func (self *Program) Pop() shared.Opcode { //nolint:ireturn
	if self.IsEmpty() {
		return nil
	}
	result := (*self)[0]
	*self = (*self)[1:]
	return result
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

func (self Program) RemoveAllButLast(removeType string) Program {
	allIndexes := slice.FindAll(self.OpcodeTypes(), removeType)
	indexesToRemove := slice.TruncateLast(allIndexes)
	return slice.RemoveAt(self, indexesToRemove...)
}

// RemoveDuplicateCheckout removes checkout opcodes that immediately follow each other from this program.
func (self *Program) RemoveDuplicateCheckout() {
	result := make([]shared.Opcode, 0, len(*self))
	// this one is populated only if the last opcode is a checkout
	var lastOpcode shared.Opcode
	for _, opcode := range *self {
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
	*self = result
}

// Implementation of the fmt.Stringer interface.
func (self Program) String() string {
	return self.StringIndented("")
}

func (self Program) StringIndented(indent string) string {
	sb := strings.Builder{}
	if self.IsEmpty() {
		sb.WriteString("(empty program)\n")
	} else {
		sb.WriteString("Program:\n")
		for o, opcode := range self {
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
		*self = make([]shared.Opcode, len(jsonOpcodes))
		for j, jsonOpcode := range jsonOpcodes {
			(*self)[j] = jsonOpcode.Opcode
		}
	}
	return nil
}
