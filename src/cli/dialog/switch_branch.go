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
			m.MoveCursorUp()
		case tea.KeyDown, tea.KeyTab:
			m.MoveCursorDown()
		case tea.KeyEnter:
			m.SelectedBranch = m.Branches[m.cursor]
			return m, tea.Quit
		case tea.KeyCtrlC:
			m.SelectedBranch = m.InitialBranch
			return m, tea.Quit
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "k":
				m.MoveCursorUp()
			case "j":
				m.MoveCursorDown()
			case "o":
				m.SelectedBranch = m.Branches[m.cursor]
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *Model) MoveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	} else {
		m.cursor = len(m.Branches) - 1
	}
}

func (m *Model) MoveCursorDown() {
	if m.cursor < len(m.Branches)-1 {
		m.cursor++
	} else {
		m.cursor = 0
	}
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
