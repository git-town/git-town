package dialogcomponents

import (
	"fmt"
	"sort"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Inputs contains the input for all dialogs in an end-to-end test.
// This struct is always mutable, so doesn't need to be wrapped in Mutable.
type Inputs struct {
	cursor Mutable[int] // index of the input to return next
	inputs []Input      // the input values
	len    int          // the total number of inputs
}

func (self Inputs) IsEmpty() bool {
	return self.cursor.Immutable() == self.len
}

// Next provides the Input for the next dialog in an end-to-end test.
func (self Inputs) Next() Option[Input] {
	if self.len == 0 {
		return None[Input]()
	}
	if self.cursor.Immutable() == self.len {
		panic("not enough dialog inputs")
	}
	result := self.inputs[self.cursor.Immutable()]
	*self.cursor.Value += 1
	return Some(result)
}

func (self Inputs) VerifyAllUsed() {
	if !self.IsEmpty() {
		panic("unused dialog inputs")
	}
}

// LoadInputs provides the Inputs to use in an end-to-end test,
// taken from the given environment variable snapshot.
func LoadInputs(environmenttVariables []string) Inputs {
	inputs := []Input{}
	sort.Strings(environmenttVariables)
	for _, environmentVariable := range environmenttVariables {
		if !strings.HasPrefix(environmentVariable, InputKey) {
			continue
		}
		_, value, match := strings.Cut(environmentVariable, "=")
		if !match {
			fmt.Printf(messages.SettingIgnoreInvalid, environmentVariable)
			continue
		}
		input := ParseInput(value)
		inputs = append(inputs, input)
	}
	return NewInputs(inputs...)
}

func NewInputs(inputs ...Input) Inputs {
	cursor := 0
	return Inputs{
		cursor: NewMutable(&cursor),
		inputs: inputs,
		len:    len(inputs),
	}
}
