package dialog

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/muesli/termenv"
)

// EnterMainBranch lets the user select a new main branch for this repo.
// This includes asking the user and updating the respective setting.
func EnterMainBranch(localBranches gitdomain.LocalBranchNames, oldMainBranch gitdomain.LocalBranchName) (selectedBranch gitdomain.LocalBranchName, abort bool, err error) {
	dialogData := mainBranchModel{}
	dialogProcess := tea.NewProgram(dialogData, tea.WithOutput(os.Stderr))
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", false, err
	}
	result := dialogResult.(mainBranchModel) //nolint:forcetypeassert
	selectedBranchName := result.selectedEntry()
	selectedBranch = gitdomain.LocalBranchName(selectedBranchName)
	return selectedBranch, result.abort, nil
}

type mainBranchModel struct {
	entries        []string
	colors         dialogColors
	cursor         int
	abort          bool
	selectionColor termenv.Style // color for the currently selected entry
}

func (self mainBranchModel) Init() tea.Cmd {
	return nil
}

func (self mainBranchModel) moveCursorDown() mainBranchModel {
	if self.cursor < len(self.entries)-1 {
		self.cursor++
	} else {
		self.cursor = 0
	}
	return self
}

func (self mainBranchModel) moveCursorUp() mainBranchModel {
	if self.cursor > 0 {
		self.cursor--
	} else {
		self.cursor = len(self.entries) - 1
	}
	return self
}

func (self mainBranchModel) selectedEntry() string {
	return self.entries[self.cursor]
}

func (self mainBranchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		self.abort = true
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
			self.abort = true
			return self, tea.Quit
		}
	}
	return self, nil
}

func (self mainBranchModel) View() string {
	s := strings.Builder{}
	for i, branch := range self.entries {
		if i == self.cursor {
			s.WriteString(self.selectionColor.Styled("> " + branch))
		} else {
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

func mainBranchPrompt(mainBranch gitdomain.LocalBranchName) string {
	result := "Please specify the main development branch:"
	if !mainBranch.IsEmpty() {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(mainBranch.String())
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}
