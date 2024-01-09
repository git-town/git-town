package dialog

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// EnterMainBranch lets the user select a new main branch for this repo.
// This includes asking the user and updating the respective setting.
func EnterMainBranch(localBranches gitdomain.LocalBranchNames, oldMainBranch gitdomain.LocalBranchName) (selectedBranch gitdomain.LocalBranchName, aborted bool, err error) {
	dialogData := mainBranchModel{
		bubbleList: newBubbleList(localBranches.Strings(), oldMainBranch.String()),
	}
	dialogResult, err := tea.NewProgram(dialogData).Run()
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	result := dialogResult.(mainBranchModel) //nolint:forcetypeassert
	selectedBranch = gitdomain.LocalBranchName(result.selectedEntry())
	return selectedBranch, result.aborted, nil
}

type mainBranchModel struct {
	bubbleList
}

func (self mainBranchModel) Init() tea.Cmd {
	return nil
}

func (self mainBranchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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

func (self mainBranchModel) View() string {
	s := strings.Builder{}
	s.WriteString("Let's start by configuring the main development branch.\n")
	s.WriteString("This is the branch from which you cut new feature branches,\n")
	s.WriteString("and into which you ship feature branches when they are done.\n")
	s.WriteString("In most repositories, this is the \"main\", \"master\", or \"development\" branch.\n\n")
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
