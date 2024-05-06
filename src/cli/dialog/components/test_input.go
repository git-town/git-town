package components

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v14/src/messages"
)

// TestInputKey specifies the name of environment variables containing input for dialogs in end-to-end tests.
const TestInputKey = "GITTOWN_DIALOG_INPUT"

// TestInput contains the input for a single dialog in an end-to-end test.
type TestInput []tea.Msg

// TestInputs contains the input for all dialogs in an end-to-end test.
type TestInputs struct {
	inputs *[]TestInput
}

// Next provides the TestInput for the next dialog in an end-to-end test.
func (self *TestInputs) Append(input TestInput) {
	*self.inputs = append(*self.inputs, input)
}

// Next provides the TestInput for the next dialog in an end-to-end test.
func (self *TestInputs) Next() TestInput {
	if len(*self.inputs) == 0 {
		return TestInput{}
	}
	result := (*self.inputs)[0]
	*self.inputs = (*self.inputs)[1:]
	return result
}

// LoadTestInputs provides the TestInputs to use in an end-to-end test,
// taken from the given environment variable snapshot.
func LoadTestInputs(environmenttVariables []string) TestInputs {
	result := NewTestInputs()
	sort.Strings(environmenttVariables)
	for _, environmentVariable := range environmenttVariables {
		if !strings.HasPrefix(environmentVariable, TestInputKey) {
			continue
		}
		_, value, match := strings.Cut(environmentVariable, "=")
		if !match {
			fmt.Printf(messages.SettingIgnoreInvalid, environmentVariable)
			continue
		}
		input := ParseTestInput(value)
		result.Append(input)
	}
	return result
}

func NewTestInputs(inputs ...TestInput) TestInputs {
	return TestInputs{
		inputs: &inputs,
	}
}

// ParseTestInput converts the given input data in the environment variable format
// into the format understood by Git Town's dialogs.
func ParseTestInput(envData string) TestInput {
	result := TestInput{}
	for _, input := range strings.Split(envData, "|") {
		if len(input) > 0 {
			result = append(result, recognizeTestInput(input))
		}
	}
	return result
}

// recognizeTestInput provides the matching BubbleTea message for the given string.
func recognizeTestInput(input string) tea.Msg { //nolint:ireturn
	switch input {
	case "backspace":
		return tea.KeyMsg{Type: tea.KeyBackspace} //exhaustruct:ignore
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC} //exhaustruct:ignore
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown} //exhaustruct:ignore
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter} //exhaustruct:ignore
	case "space":
		return tea.KeyMsg{Type: tea.KeySpace} //exhaustruct:ignore
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp} //exhaustruct:ignore
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc} //exhaustruct:ignore
	case "0":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'0'}} //exhaustruct:ignore
	case "1":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}} //exhaustruct:ignore
	case "2":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}} //exhaustruct:ignore
	case "3":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}} //exhaustruct:ignore
	case "4":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}} //exhaustruct:ignore
	case "5":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}} //exhaustruct:ignore
	case "6":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'6'}} //exhaustruct:ignore
	case "7":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'7'}} //exhaustruct:ignore
	case "8":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'8'}} //exhaustruct:ignore
	case "9":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}} //exhaustruct:ignore
	case "a":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}} //exhaustruct:ignore
	case "c":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}} //exhaustruct:ignore
	case "d":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}} //exhaustruct:ignore
	case "e":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}} //exhaustruct:ignore
	case "n":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}} //exhaustruct:ignore
	case "o":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}} //exhaustruct:ignore
	case "q":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}} //exhaustruct:ignore
	}
	panic("unknown test input: " + input)
}
