package dialog

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const TestInputKey = "GITTOWN_DIALOG_INPUT"

type TestInput []tea.Msg

type TestInputs []TestInput

func (self *TestInputs) Next() TestInput {
	if len(*self) == 0 {
		return TestInput{}
	}
	result := (*self)[0]
	*self = (*self)[1:]
	return result
}

func LoadTestInputs(environmenttVariables []string) TestInputs {
	result := TestInputs{}
	sort.Strings(environmenttVariables)
	for _, environmentVariable := range environmenttVariables {
		if !strings.HasPrefix(environmentVariable, TestInputKey) {
			continue
		}
		_, value, match := strings.Cut(environmentVariable, "=")
		if !match {
			fmt.Printf("Notice: ignoring invalid dialog input setting %q\n", environmentVariable)
			continue
		}
		inputs := ParseTestInput(value)
		result = append(result, inputs)
	}
	return result
}

func ParseTestInput(envData string) TestInput {
	result := TestInput{}
	for _, input := range strings.Split(envData, "|") {
		if len(input) > 0 {
			result = append(result, RecognizeTestInput(input))
		}
	}
	return result
}

func RecognizeTestInput(input string) tea.Msg { //nolint:ireturn
	switch input { //nolint:ireturn
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
