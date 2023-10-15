package program

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// JSON is used to store a step in JSON.
type JSON struct { //nolint:musttag // JSONStep uses a custom serialization algorithm
	Step shared.Opcode
}

// MarshalJSON marshals the step to JSON.
func (js *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": js.Step,
		"type": gohacks.TypeName(js.Step),
	})
}

// UnmarshalJSON unmarshals the step from JSON.
func (js *JSON) UnmarshalJSON(b []byte) error {
	var mapping map[string]json.RawMessage
	err := json.Unmarshal(b, &mapping)
	if err != nil {
		return err
	}
	var stepType string
	err = json.Unmarshal(mapping["type"], &stepType)
	if err != nil {
		return err
	}
	js.Step = opcode.DetermineStep(stepType)
	if js.Step == nil {
		return fmt.Errorf(messages.RunstateStepUnknown, stepType)
	}
	return json.Unmarshal(mapping["data"], &js.Step)
}
