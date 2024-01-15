package dialog

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

func SwitchBranch(localBranches gitdomain.LocalBranchNames, initialBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchName, error) {
	localBranchNames := localBranches.Strings()
	dialogData := switchModel{
		bubbleList:    newBubbleList(localBranchNames, DetermineCursorPos(localBranchNames, initialBranch.String())),
		initialBranch: initialBranch.String(),
	}
	dialogProcess := tea.NewProgram(dialogData, tea.WithOutput(os.Stderr))
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", err
	}
	result := dialogResult.(switchModel) //nolint:forcetypeassert
	selectedBranch := gitdomain.NewLocalBranchName(result.bubbleList.selectedEntry())
	return selectedBranch, nil
}

type switchModel struct {
	bubbleList
	initialBranch string // name of the currently checked out branch
}

func (self switchModel) Init() tea.Cmd {
	return nil
}

func (self switchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, code := self.bubbleList.handleKey(keyMsg); handled {
		return self, code
	}
	if keyMsg.Type == tea.KeyEnter {
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		return self, tea.Quit
	}
	return self, nil
}

func (self switchModel) View() string {
	s := strings.Builder{}
	for i, branch := range self.entries {
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
