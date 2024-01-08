package dialog

import (
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func SwitchBranch(branchNames []string, initialBranch string) (string, error) {
	cursor := slices.Index(branchNames, initialBranch)
	if cursor < 0 {
		cursor = 0
	}
	dialogData := SwitchModel{
		branches:       branchNames,
		cursor:         cursor,
		helpColor:      termenv.String().Faint(),
		helpKeyColor:   termenv.String().Faint().Bold(),
		initialBranch:  initialBranch,
		initialColor:   termenv.String().Foreground(termenv.ANSIGreen),
		SelectedBranch: initialBranch,
		selectionColor: termenv.String().Foreground(termenv.ANSICyan),
	}
	dialogProcess := tea.NewProgram(dialogData, tea.WithOutput(os.Stderr))
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", err
	}
	result := dialogResult.(SwitchModel) //nolint:forcetypeassert
	return result.SelectedBranch, nil
}

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

func (self SwitchModel) Init() tea.Cmd {
	return nil
}

func (self SwitchModel) MoveCursorDown() SwitchModel {
	if self.cursor < len(self.branches)-1 {
		self.cursor++
	} else {
		self.cursor = 0
	}
	return self
}

func (self SwitchModel) MoveCursorUp() SwitchModel {
	if self.cursor > 0 {
		self.cursor--
	} else {
		self.cursor = len(self.branches) - 1
	}
	return self
}

func (self SwitchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	switch keyMsg.Type { //nolint:exhaustive
	case tea.KeyUp, tea.KeyShiftTab:
		return self.MoveCursorUp(), nil
	case tea.KeyDown, tea.KeyTab:
		return self.MoveCursorDown(), nil
	case tea.KeyEnter:
		self.SelectedBranch = self.branches[self.cursor]
		return self, tea.Quit
	case tea.KeyCtrlC:
		self.SelectedBranch = self.initialBranch
		return self, tea.Quit
	case tea.KeyRunes:
		switch string(keyMsg.Runes) {
		case "k", "A", "Z":
			return self.MoveCursorUp(), nil
		case "j", "B":
			return self.MoveCursorDown(), nil
		case "o":
			self.SelectedBranch = self.branches[self.cursor]
			return self, tea.Quit
		case "q":
			self.SelectedBranch = self.initialBranch
			return self, tea.Quit
		}
	}
	return self, nil
}

func (self SwitchModel) View() string {
	s := strings.Builder{}
	for i, branch := range self.branches {
		switch {
		case i == self.cursor:
			s.WriteString(self.selectionColor.Styled("> " + branch))
		case branch == self.initialBranch:
			s.WriteString(self.initialColor.Styled("* " + branch))
		default:
			s.WriteString("  " + branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(self.helpKeyColor.Styled("↑"))
	s.WriteString(self.helpColor.Styled("/"))
	s.WriteString(self.helpKeyColor.Styled("k"))
	s.WriteString(self.helpColor.Styled(" up   "))
	// down
	s.WriteString(self.helpKeyColor.Styled("↓"))
	s.WriteString(self.helpColor.Styled("/"))
	s.WriteString(self.helpKeyColor.Styled("j"))
	s.WriteString(self.helpColor.Styled(" down   "))
	// accept
	s.WriteString(self.helpKeyColor.Styled("enter"))
	s.WriteString(self.helpColor.Styled("/"))
	s.WriteString(self.helpKeyColor.Styled("o"))
	s.WriteString(self.helpColor.Styled(" accept   "))
	// abort
	s.WriteString(self.helpKeyColor.Styled("esc"))
	s.WriteString(self.helpColor.Styled("/"))
	s.WriteString(self.helpKeyColor.Styled("q"))
	s.WriteString(self.helpColor.Styled(" abort"))
	return s.String()
}
