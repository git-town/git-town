package dialogcomponents

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// InputKey specifies the name of environment variables containing input for dialogs in end-to-end tests.
const InputKey = "GITTOWN_DIALOG_INPUT"

// Input contains the input for a single dialog in an end-to-end test.
type Input struct {
	Messages []tea.Msg
	StepName string
}

// ParseInput converts the given input data in the environment variable format
// into the format understood by Git Town's dialogs.
func ParseInput(envData string) Input {
	messages := []tea.Msg{}
	stepName, keys, has := strings.Cut(envData, "@")
	if !has {
		panic(fmt.Sprintf("found test input without step name: %q", envData))
	}
	for input := range strings.SplitSeq(keys, "|") {
		if len(input) > 0 {
			messages = append(messages, recognizeInput(input))
		}
	}
	return Input{
		Messages: messages,
		StepName: stepName,
	}
}

// recognizeInput provides the matching BubbleTea message for the given string.
func recognizeInput(input string) tea.Msg { //nolint:ireturn
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
		return tea.KeyMsg{Runes: []rune{'0'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "1":
		return tea.KeyMsg{Runes: []rune{'1'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "2":
		return tea.KeyMsg{Runes: []rune{'2'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "3":
		return tea.KeyMsg{Runes: []rune{'3'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "4":
		return tea.KeyMsg{Runes: []rune{'4'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "5":
		return tea.KeyMsg{Runes: []rune{'5'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "6":
		return tea.KeyMsg{Runes: []rune{'6'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "7":
		return tea.KeyMsg{Runes: []rune{'7'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "8":
		return tea.KeyMsg{Runes: []rune{'8'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "9":
		return tea.KeyMsg{Runes: []rune{'9'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "a":
		return tea.KeyMsg{Runes: []rune{'a'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "b":
		return tea.KeyMsg{Runes: []rune{'b'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "c":
		return tea.KeyMsg{Runes: []rune{'c'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "d":
		return tea.KeyMsg{Runes: []rune{'d'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "e":
		return tea.KeyMsg{Runes: []rune{'e'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "f":
		return tea.KeyMsg{Runes: []rune{'f'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "g":
		return tea.KeyMsg{Runes: []rune{'g'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "h":
		return tea.KeyMsg{Runes: []rune{'h'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "i":
		return tea.KeyMsg{Runes: []rune{'i'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "j":
		return tea.KeyMsg{Runes: []rune{'j'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "k":
		return tea.KeyMsg{Runes: []rune{'k'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "l":
		return tea.KeyMsg{Runes: []rune{'l'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "m":
		return tea.KeyMsg{Runes: []rune{'m'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "n":
		return tea.KeyMsg{Runes: []rune{'n'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "o":
		return tea.KeyMsg{Runes: []rune{'o'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "p":
		return tea.KeyMsg{Runes: []rune{'p'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "q":
		return tea.KeyMsg{Runes: []rune{'q'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "r":
		return tea.KeyMsg{Runes: []rune{'r'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "s":
		return tea.KeyMsg{Runes: []rune{'s'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "t":
		return tea.KeyMsg{Runes: []rune{'t'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "u":
		return tea.KeyMsg{Runes: []rune{'u'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "v":
		return tea.KeyMsg{Runes: []rune{'v'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "w":
		return tea.KeyMsg{Runes: []rune{'w'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "x":
		return tea.KeyMsg{Runes: []rune{'x'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "y":
		return tea.KeyMsg{Runes: []rune{'y'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "z":
		return tea.KeyMsg{Runes: []rune{'z'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "^":
		return tea.KeyMsg{Runes: []rune{'^'}, Type: tea.KeyRunes} //exhaustruct:ignore
	case "-":
		return tea.KeyMsg{Runes: []rune{'-'}, Type: tea.KeyRunes} //exhaustruct:ignore
	}
	panic("unknown test input: " + input)
}
