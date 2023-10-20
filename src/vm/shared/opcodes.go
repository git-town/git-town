package shared

import "encoding/json"

type Opcodes []Opcode

// MarshalJSON marshals this program to JSON.
func (self *Opcodes) MarshalJSON() ([]byte, error) {
	jsonOpcodes := make([]JSON, len(self.Opcodes))
	for o, opcode := range self.Opcodes {
		jsonOpcodes[o] = JSON{Opcode: opcode}
	}
	return json.Marshal(jsonOpcodes)
}
