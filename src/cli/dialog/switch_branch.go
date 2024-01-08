package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type SwitchModel struct {
	activeColor    *color.Color
	branches       []string // names of all branches
	cursor         int      // 0-based number of the selected row
	initialBranch  string   // name of the currently checked out branch
	initialColor   *color.Color
	SelectedBranch string // name of the currently selected branch
}

func NewSwitchModel(branches []string, initialBranch string) SwitchModel {
	cursor := slices.Index(branches, initialBranch)
	if cursor == -1 {
		cursor = 0
	}
	return SwitchModel{
		activeColor:    color.New(color.FgCyan),
		branches:       branches,
		cursor:         cursor,
		initialBranch:  initialBranch,
		initialColor:   color.New(color.FgGreen),
		SelectedBranch: initialBranch,
	}
}

func (m SwitchModel) Init() tea.Cmd {
	return nil
}

func (m SwitchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp, tea.KeyShiftTab:
			return m.MoveCursorUp(), nil
		case tea.KeyDown, tea.KeyTab:
			return m.MoveCursorDown(), nil
		case tea.KeyEnter:
			m.SelectedBranch = m.branches[m.cursor]
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.SelectedBranch = m.initialBranch
			return m, tea.Quit
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "k":
				return m.MoveCursorUp(), nil
			case "j":
				return m.MoveCursorDown(), nil
			case "o":
				m.SelectedBranch = m.branches[m.cursor]
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m SwitchModel) MoveCursorUp() SwitchModel {
	if m.cursor > 0 {
		m.cursor--
	} else {
		m.cursor = len(m.branches) - 1
	}
	return m
}

func (m SwitchModel) MoveCursorDown() SwitchModel {
	if m.cursor < len(m.branches)-1 {
		m.cursor++
	} else {
		m.cursor = 0
	}
	return m
}

func (m SwitchModel) View() string {
	s := strings.Builder{}
	for i, branch := range m.branches {
		if i == m.cursor {
			s.WriteString(m.activeColor.Sprint("> "))
			s.WriteString(m.activeColor.Sprint(branch))
		} else if branch == m.initialBranch {
			s.WriteString(m.initialColor.Sprint("* "))
			s.WriteString(m.initialColor.Sprint(branch))
		} else {
			s.WriteString("  ")
			s.WriteString(branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n")
	s.WriteString("Press Ctrl-C to quit")
	return s.String()
}
