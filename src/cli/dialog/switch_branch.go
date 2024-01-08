package dialog

import (
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/muesli/termenv"
)

type SwitchModel struct {
	activeColor        *color.Color
	branches           []string // names of all branches
	cursor             int      // 0-based number of the selected row
	helpColor          *color.Color
	helpHighlightColor *color.Color
	initialBranch      string // name of the currently checked out branch
	initialColor       *color.Color
	SelectedBranch     string // name of the currently selected branch
}

func NewSwitchModel(branches []string, initialBranch string) SwitchModel {
	cursor := slices.Index(branches, initialBranch)
	if cursor == -1 {
		cursor = 0
	}
	output := termenv.NewOutput(os.Stderr)
	darkTheme := output.HasDarkBackground()
	fmt.Println(termenv.String("11111111111111").Foreground(termenv.ANSIBlue))
	s := output.String("Hello World")
	fmt.Println("DARK THEME:", darkTheme)
	fmt.Println("FOREGROUND:", output.ForegroundColor())
	fmt.Println("BACKGROUND:", output.BackgroundColor())
	s.Background(output.Color("3"))
	s.Bold()
	fmt.Println(s)

	return SwitchModel{
		activeColor:        color.New(color.FgCyan),
		branches:           branches,
		cursor:             cursor,
		helpColor:          color.New(color.Faint),
		helpHighlightColor: color.New(color.Faint).Add(color.Bold),
		initialBranch:      initialBranch,
		initialColor:       color.New(color.FgGreen),
		SelectedBranch:     initialBranch,
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
	if isKeyMsg {
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
	}
	return m, nil
}

func (m SwitchModel) View() string {
	s := strings.Builder{}
	for i, branch := range m.branches {
		switch {
		case i == m.cursor:
			s.WriteString(m.activeColor.Sprint("> "))
			s.WriteString(m.activeColor.Sprint(branch))
		case branch == m.initialBranch:
			s.WriteString(m.initialColor.Sprint("* "))
			s.WriteString(m.initialColor.Sprint(branch))
		default:
			s.WriteString("  ")
			s.WriteString(branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n")
	s.WriteString("  ")
	// up
	s.WriteString(termenv.String("↑").Faint().Bold().String())
	s.WriteString(m.helpColor.Sprint("/"))
	s.WriteString(m.helpHighlightColor.Sprint("k"))
	s.WriteString(m.helpColor.Sprint(" up   "))
	// down
	s.WriteString(m.helpHighlightColor.Sprint("↓"))
	s.WriteString(m.helpColor.Sprint("/"))
	s.WriteString(m.helpHighlightColor.Sprint("j"))
	s.WriteString(m.helpColor.Sprint(" down   "))
	// accept
	s.WriteString(m.helpHighlightColor.Sprint("enter"))
	s.WriteString(m.helpColor.Sprint("/"))
	s.WriteString(m.helpHighlightColor.Sprint("o"))
	s.WriteString(m.helpColor.Sprint(" accept   "))
	// abort
	s.WriteString(m.helpHighlightColor.Sprint("esc"))
	s.WriteString(m.helpColor.Sprint("/"))
	s.WriteString(m.helpHighlightColor.Sprint("q"))
	s.WriteString(m.helpColor.Sprint(" abort"))
	return s.String()
}
