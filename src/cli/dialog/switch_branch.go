package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type SwitchModel struct {
	activeColor        *color.Color
	branches           []string // names of all branches
	cursor             int      // 0-based number of the selected row
	helpColor          *color.Color
	helpHighlightColor *color.Color
	initialBranch      string // name of the currently checked out branch
	initialColor       *color.Color
	showHelp           bool
	SelectedBranch     string // name of the currently selected branch
}

func NewSwitchModel(branches []string, initialBranch string) SwitchModel {
	cursor := slices.Index(branches, initialBranch)
	if cursor == -1 {
		cursor = 0
	}
	return SwitchModel{
		activeColor:        color.New(color.FgCyan),
		branches:           branches,
		cursor:             cursor,
		helpColor:          color.New(color.Faint),
		helpHighlightColor: color.New(color.Faint).Add(color.Bold),
		initialBranch:      initialBranch,
		initialColor:       color.New(color.FgGreen),
		SelectedBranch:     initialBranch,
		showHelp:           false,
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
			case "h", "?":
				m.showHelp = !m.showHelp
				return m, nil
			case "k":
				return m.MoveCursorUp(), nil
			case "j":
				return m.MoveCursorDown(), nil
			case "o":
				m.SelectedBranch = m.branches[m.cursor]
				return m, tea.Quit
			case "q":
				m.SelectedBranch = m.initialBranch
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
	if m.showHelp {
		s.WriteString(m.helpColor.Sprint("[down] or \"j\"             select the next branch\n"))
		s.WriteString(m.helpHighlightColor.Sprint("[up] or \"k\"               select the previous branch\n"))
		s.WriteString(m.helpColor.Sprint("[enter] or \"o\"            finish and check out the selected branch\n"))
		s.WriteString(m.helpColor.Sprint("[esc] or [ctrl-c] or \"q\"  quit without changing the branch\n"))
		s.WriteString(m.helpColor.Sprint("\"h\"                       toggle this help message\n"))
	} else {
		s.WriteString("  ")
		s.WriteString(m.helpHighlightColor.Sprint("↑"))
		s.WriteString(m.helpColor.Sprint("/"))
		s.WriteString(m.helpHighlightColor.Sprint("k"))
		s.WriteString(m.helpColor.Sprint(" up   "))
		s.WriteString(m.helpHighlightColor.Sprint("↓"))
		s.WriteString(m.helpColor.Sprint("/"))
		s.WriteString(m.helpHighlightColor.Sprint("j"))
		s.WriteString(m.helpColor.Sprint(" down  "))

	}
	return s.String()
}
