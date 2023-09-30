package runstate

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/steps"
)

// JSONStep is used to store a step in JSON.
type JSONStep struct { //nolint:musttag // JSONStep uses a custom serialization algorithm
	Step steps.Step
}

// MarshalJSON marshals the step to JSON.
func (j *JSONStep) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": j.Step,
		"type": gohacks.TypeName(j.Step),
	})
}

// UnmarshalJSON unmarshals the step from JSON.
func (j *JSONStep) UnmarshalJSON(b []byte) error {
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
	j.Step = steps.All[stepType]
	if j.Step == nil {
		return fmt.Errorf(messages.RunstateStepUnknown, stepType)
	}
	return json.Unmarshal(mapping["data"], &j.Step)
}
