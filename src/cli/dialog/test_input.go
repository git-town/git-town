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
	case "enter":
		return tea.KeyEnter
	case "0":
		return tea.Key{Type: tea.KeyRunes, Runes: []rune{'0'}, Alt: false}
	}
	panic("unknown test input: " + input)
}
