package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// EnterMainBranch lets the user select a new main branch for this repo.
// This includes asking the user and updating the respective setting.
func EnterMainBranch(localBranches gitdomain.LocalBranchNames, oldMainBranch gitdomain.LocalBranchName) (selectedBranch gitdomain.LocalBranchName, aborted bool, err error) {
	cursor := slices.Index(localBranches, oldMainBranch)
	if cursor < 0 {
		cursor = 0
	}
	dialogData := mainBranchModel{
		aborted: false,
		entries: localBranches.Strings(),
		colors:  createColors(),
		cursor:  cursor,
	}
	dialogProcess := tea.NewProgram(dialogData)
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	result := dialogResult.(mainBranchModel) //nolint:forcetypeassert
	selectedBranch = gitdomain.LocalBranchName(result.selectedEntry())
	return selectedBranch, result.aborted, nil
}

type mainBranchModel struct {
	aborted bool
	colors  dialogColors
	cursor  int
	entries []string
}

func (self mainBranchModel) Init() tea.Cmd {
	return nil
}

func (self mainBranchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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
		self.aborted = true
		return self, tea.Quit
	}
	switch keyMsg.String() {
	case "k", "A", "Z":
		return self.moveCursorUp(), nil
	case "j", "B":
		return self.moveCursorDown(), nil
	case "o":
		return self, tea.Quit
	case "q":
		self.aborted = true
		return self, tea.Quit
	}
	return self, nil
}

func (self mainBranchModel) View() string {
	s := strings.Builder{}
	s.WriteString("Let's start by setting up the main development branch.\n")
	s.WriteString("This is the branch from which you mostly cut new feature branches,\n")
	s.WriteString("and into which you ship feature branches when they are done.\n")
	s.WriteString("In most repositories, this branch is called \"main\", \"master\", or \"development\".\n\n")
	for i, branch := range self.entries {
		if i == self.cursor {
			s.WriteString(self.colors.selection.Styled("> " + branch))
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
