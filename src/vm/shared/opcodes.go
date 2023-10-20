package shared

import (
	"encoding/json"

	"github.com/git-town/git-town/v9/src/vm/program"
)

type Opcodes []Opcode

// MarshalJSON marshals this program to JSON.
func (self Opcodes) MarshalJSON() ([]byte, error) {
	jsonOpcodes := make([]program.JSON, len(self))
	for o, opcode := range self {
		jsonOpcodes[o] = program.JSON{Opcode: opcode}
	}
	return json.Marshal(jsonOpcodes)
}
