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

func RecognizeTestInput(input string) tea.Msg {
	switch input {
	case "enter":
		return tea.KeyEnter
	}
	panic("unknown test input: " + input)
}
