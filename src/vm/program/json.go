package program

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/persistence"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

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
	self.Opcode = persistence.DetermineOpcode(opcodeType)
	if self.Opcode == nil {
		return fmt.Errorf(messages.OpcodeUnknown, opcodeType)
	}
	return json.Unmarshal(mapping["data"], &self.Opcode)
}
