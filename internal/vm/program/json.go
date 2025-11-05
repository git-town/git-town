package program

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// JSON is used to store an opcode in JSON.
type JSON struct {
	shared.Opcode
}

// MarshalJSON marshals the opcode to JSON.
func (self *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"data": self.Opcode,
		"type": gohacks.TypeName(self.Opcode),
	})
}

// UnmarshalJSON unmarshals the opcode from JSON.
func (self *JSON) UnmarshalJSON(b []byte) error {
	var mapping map[string]json.RawMessage
	if err := json.Unmarshal(b, &mapping); err != nil {
		return err
	}
	var opcodeType string
	if err := json.Unmarshal(mapping["type"], &opcodeType); err != nil {
		return err
	}
	self.Opcode = opcodes.Lookup(opcodeType)
	if self.Opcode == nil {
		return fmt.Errorf(messages.OpcodeUnknown, opcodeType)
	}
	return json.Unmarshal(mapping["data"], &self.Opcode)
}
