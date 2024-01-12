package dialog

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const TestInputKey = "GITTOWN_DIALOG_INPUT"

func ParseTestInput(envData string) []tea.Msg {
	result := []tea.Msg{}
	for _, input := range strings.Split(envData, "|") {
		result = append(result, RecognizeTestInput(input))
	}
	return result
}

func RecognizeTestInput(input string) tea.Msg { //nolint:ireturn
	switch input {
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC} //nolint:exhaustruct
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown} //nolint:exhaustruct
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter} //nolint:exhaustruct
	case "o":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}} //nolint:exhaustruct
	case "q":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}} //nolint:exhaustruct
	case "space":
		return tea.KeyMsg{Type: tea.KeySpace} //nolint:exhaustruct
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp} //nolint:exhaustruct
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
	}
	panic("unknown test input: " + input)
}
