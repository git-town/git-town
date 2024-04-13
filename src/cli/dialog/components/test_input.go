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
type TestInputs []TestInput

// Next provides the TestInput for the next dialog in an end-to-end test.
func (self *TestInputs) Next() TestInput {
	if len(*self) == 0 {
		return TestInput{}
	}
	result := (*self)[0]
	*self = (*self)[1:]
	return result
}

// LoadTestInputs provides the TestInputs to use in an end-to-end test,
// taken from the given environment variable snapshot.
func LoadTestInputs(environmenttVariables []string) TestInputs {
	result := TestInputs{}
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
		inputs := ParseTestInput(value)
		result = append(result, inputs)
	}
	return result
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
		return tea.KeyMsg{Type: tea.KeyBackspace} //nolint:exhaustruct
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC} //nolint:exhaustruct
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown} //nolint:exhaustruct
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter} //nolint:exhaustruct
	case "space":
		return tea.KeyMsg{Type: tea.KeySpace} //nolint:exhaustruct
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp} //nolint:exhaustruct
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc} //nolint:exhaustruct
	case "0":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'0'}} //nolint:exhaustruct
	case "1":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}} //nolint:exhaustruct
	case "2":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}} //nolint:exhaustruct
	case "3":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}} //nolint:exhaustruct
	case "4":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}} //nolint:exhaustruct
	case "5":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}} //nolint:exhaustruct
	case "6":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'6'}} //nolint:exhaustruct
	case "7":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'7'}} //nolint:exhaustruct
	case "8":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'8'}} //nolint:exhaustruct
	case "9":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}} //nolint:exhaustruct
	case "a":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}} //nolint:exhaustruct
	case "c":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}} //nolint:exhaustruct
	case "d":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}} //nolint:exhaustruct
	case "e":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}} //nolint:exhaustruct
	case "n":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}} //nolint:exhaustruct
	case "o":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}} //nolint:exhaustruct
	case "q":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}} //nolint:exhaustruct
	}
	panic("unknown test input: " + input)
}
