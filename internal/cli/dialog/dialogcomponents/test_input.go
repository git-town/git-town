package dialogcomponents

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// TestInputKey specifies the name of environment variables containing input for dialogs in end-to-end tests.
const TestInputKey = "GITTOWN_DIALOG_INPUT"

// TestInput contains the input for a single dialog in an end-to-end test.
type TestInput struct {
	Messages []tea.Msg
	StepName string
}

// ParseTestInput converts the given input data in the environment variable format
// into the format understood by Git Town's dialogs.
func ParseTestInput(envData string) TestInput {
	messages := []tea.Msg{}
	stepName, keys, has := strings.Cut(envData, "@")
	if !has {
		panic(fmt.Sprintf("found test input without step name: %q", envData))
	}
	for _, input := range strings.Split(keys, "|") {
		if len(input) > 0 {
			messages = append(messages, recognizeTestInput(input))
		}
	}
	return TestInput{messages, stepName}
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
	case "k":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}} //exhaustruct:ignore
	case "n":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}} //exhaustruct:ignore
	case "o":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}} //exhaustruct:ignore
	case "p":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}} //exhaustruct:ignore
	case "q":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}} //exhaustruct:ignore
	case "r":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}} //exhaustruct:ignore
	case "s":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}} //exhaustruct:ignore
	case "u":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}} //exhaustruct:ignore
	}
	panic("unknown test input: " + input)
}
