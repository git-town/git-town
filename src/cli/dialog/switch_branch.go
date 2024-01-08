package dialog

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Branches       []string // names of all branches
	cursor         int      // 0-based number of the selected row
	InitialBranch  string   // name of the currently checked out branch
	SelectedBranch string   // name of the currently selected branch
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp, tea.KeyShiftTab:
			return m.MoveCursorUp(), nil
		case tea.KeyDown, tea.KeyTab:
			return m.MoveCursorDown(), nil
		case tea.KeyEnter:
			m.SelectedBranch = m.Branches[m.cursor]
			return m, tea.Quit
		case tea.KeyCtrlC:
			m.SelectedBranch = m.InitialBranch
			return m, tea.Quit
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "k":
				return m.MoveCursorUp(), nil
			case "j":
				return m.MoveCursorDown(), nil
			case "o":
				m.SelectedBranch = m.Branches[m.cursor]
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) MoveCursorUp() Model {
	if m.cursor > 0 {
		m.cursor--
	} else {
		m.cursor = len(m.Branches) - 1
	}
	return m
}

func (m Model) MoveCursorDown() Model {
	if m.cursor < len(m.Branches)-1 {
		m.cursor++
	} else {
		m.cursor = 0
	}
	return m
}

func (m Model) View() string {
	s := strings.Builder{}
	for _, branch := range m.Branches {
		if branch == m.InitialBranch {
			s.WriteString("> ")
			s.WriteString(branch)
		} else if branch == m.SelectedBranch {
			s.WriteString("* ")
			s.WriteString(branch)
		} else {
			s.WriteString("  ")
			s.WriteString(branch)
		}
		s.WriteRune('\n')
	}

	s.WriteString("Press Ctrl-C to quit")
	return s.String()
}
