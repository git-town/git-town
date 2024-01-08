package dialog

import (
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func SwitchBranch(branchNames []string, initialBranch string) (string, error) {
	cursor := slices.Index(branchNames, initialBranch)
	if cursor < 0 {
		cursor = 0
	}
	dialogData := switchModel{
		branches:      branchNames,
		cursor:        cursor,
		colors:        createColors(),
		initialBranch: initialBranch,
	}
	dialogProcess := tea.NewProgram(dialogData, tea.WithOutput(os.Stderr))
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", err
	}
	result := dialogResult.(switchModel) //nolint:forcetypeassert
	return result.selectedBranch(), nil
}

type switchModel struct {
	branches      []string     // names of all branches
	cursor        int          // index of the currently selected row
	colors        dialogColors // colors to use in the dialog
	initialBranch string       // name of the currently checked out branch
}

func (self switchModel) Init() tea.Cmd {
	return nil
}

func (self switchModel) moveCursorDown() switchModel {
	if self.cursor < len(self.branches)-1 {
		self.cursor++
	} else {
		self.cursor = 0
	}
	return self
}

func (self switchModel) moveCursorUp() switchModel {
	if self.cursor > 0 {
		self.cursor--
	} else {
		self.cursor = len(self.branches) - 1
	}
	return self
}

func (self switchModel) selectedBranch() string {
	return self.branches[self.cursor]
}

func (self switchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	switch keyMsg.Type { //nolint:exhaustive
	case tea.KeyUp, tea.KeyShiftTab:
		return self.moveCursorUp(), nil
	case tea.KeyDown, tea.KeyTab:
		return self.moveCursorDown(), nil
	case tea.KeyEnter:
		return self, tea.Quit
	case tea.KeyCtrlC:
		return self, tea.Quit
	case tea.KeyRunes:
		switch string(keyMsg.Runes) {
		case "k", "A", "Z":
			return self.moveCursorUp(), nil
		case "j", "B":
			return self.moveCursorDown(), nil
		case "o":
			return self, tea.Quit
		case "q":
			return self, tea.Quit
		}
	}
	return self, nil
}

func (self switchModel) View() string {
	s := strings.Builder{}
	for i, branch := range self.branches {
		switch {
		case i == self.cursor:
			s.WriteString(self.colors.selection.Styled("> " + branch))
		case branch == self.initialBranch:
			s.WriteString(self.colors.initial.Styled("* " + branch))
		default:
			s.WriteString("  " + branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(self.colors.helpKey.Styled("↑"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("k"))
	s.WriteString(self.colors.help.Styled(" up   "))
	// down
	s.WriteString(self.colors.helpKey.Styled("↓"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("j"))
	s.WriteString(self.colors.help.Styled(" down   "))
	// accept
	s.WriteString(self.colors.helpKey.Styled("enter"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("o"))
	s.WriteString(self.colors.help.Styled(" accept   "))
	// abort
	s.WriteString(self.colors.helpKey.Styled("esc"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("q"))
	s.WriteString(self.colors.help.Styled(" abort"))
	return s.String()
}
