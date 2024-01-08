package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type SwitchModel struct {
	branches       []string      // names of all branches
	cursor         int           // index of the currently selected row
	helpColor      termenv.Style // color of help text
	helpKeyColor   termenv.Style // color of key names in help text
	initialBranch  string        // name of the currently checked out branch
	initialColor   termenv.Style // color for the row containing the currently checked out branch
	SelectedBranch string        // name of the currently selected branch
	selectionColor termenv.Style // color for the currently selected entry
}

func NewSwitchModel(branches []string, initialBranch string) SwitchModel {
	cursor := slices.Index(branches, initialBranch)
	if cursor < 0 {
		cursor = 0
	}
	return SwitchModel{
		branches:       branches,
		cursor:         cursor,
		helpColor:      termenv.String().Faint(),
		helpKeyColor:   termenv.String().Faint().Bold(),
		initialBranch:  initialBranch,
		initialColor:   termenv.String().Foreground(termenv.ANSIGreen),
		SelectedBranch: initialBranch,
		selectionColor: termenv.String().Foreground(termenv.ANSICyan),
	}
}

func (m SwitchModel) Init() tea.Cmd {
	return nil
}

func (m SwitchModel) MoveCursorDown() SwitchModel {
	if m.cursor < len(m.branches)-1 {
		m.cursor++
	} else {
		m.cursor = 0
	}
	return m
}

func (m SwitchModel) MoveCursorUp() SwitchModel {
	if m.cursor > 0 {
		m.cursor--
	} else {
		m.cursor = len(m.branches) - 1
	}
	return m
}

func (m SwitchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return m, nil
	}
	switch keyMsg.Type { //nolint:exhaustive
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
		switch string(keyMsg.Runes) {
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
	return m, nil
}

func (m SwitchModel) View() string {
	s := strings.Builder{}
	for i, branch := range m.branches {
		switch {
		case i == m.cursor:
			s.WriteString(m.selectionColor.Styled("> " + branch))
		case branch == m.initialBranch:
			s.WriteString(m.initialColor.Styled("* " + branch))
		default:
			s.WriteString("  " + branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(m.helpKeyColor.Styled("↑"))
	s.WriteString(m.helpColor.Styled("/"))
	s.WriteString(m.helpKeyColor.Styled("k"))
	s.WriteString(m.helpColor.Styled(" up   "))
	// down
	s.WriteString(m.helpKeyColor.Styled("↓"))
	s.WriteString(m.helpColor.Styled("/"))
	s.WriteString(m.helpKeyColor.Styled("j"))
	s.WriteString(m.helpColor.Styled(" down   "))
	// accept
	s.WriteString(m.helpKeyColor.Styled("enter"))
	s.WriteString(m.helpColor.Styled("/"))
	s.WriteString(m.helpKeyColor.Styled("o"))
	s.WriteString(m.helpColor.Styled(" accept   "))
	// abort
	s.WriteString(m.helpKeyColor.Styled("esc"))
	s.WriteString(m.helpColor.Styled("/"))
	s.WriteString(m.helpKeyColor.Styled("q"))
	s.WriteString(m.helpColor.Styled(" abort"))
	return s.String()
}
