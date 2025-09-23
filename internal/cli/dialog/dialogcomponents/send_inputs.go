package dialogcomponents

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// SendInputs sends the given keystrokes to the given bubbletea program.
func SendInputs(stepName string, input Option[Input], program *tea.Program) {
	if input, has := input.Get(); has {
		if stepName != input.StepName {
			panic(fmt.Sprintf("mismatching dialog names: want %q but have %q", stepName, input.StepName))
		}
		go func() {
			for _, msg := range input.Messages {
				program.Send(msg)
			}
		}()
	}
}
