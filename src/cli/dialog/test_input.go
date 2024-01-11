package dialog

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "o":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}, Alt: false}
	case "q":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}, Alt: false}
	case "space":
		return tea.KeyMsg{Type: tea.KeySpace}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "0":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'0'}, Alt: false}
	case "1":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}, Alt: false}
	case "2":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}, Alt: false}
	case "3":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}, Alt: false}
	case "4":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}, Alt: false}
	case "5":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}, Alt: false}
	case "6":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'6'}, Alt: false}
	case "7":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'7'}, Alt: false}
	case "8":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'8'}, Alt: false}
	case "9":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}, Alt: false}
	}
	panic("unknown test input: " + input)
}
