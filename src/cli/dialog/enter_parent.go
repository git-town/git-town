package dialog

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

const PerennialBranchOption = "<none> (perennial branch)"

// EnterMainBranch lets the user select a new main branch for this repo.
// This includes asking the user and updating the respective setting.
func EnterParent(branch gitdomain.LocalBranchName, localBranches gitdomain.LocalBranchNames, lineage configdomain.Lineage, mainBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchName, bool, error) {
	parentCandidates := localBranches.Remove(branch).Remove(lineage.Children(branch)...).Strings()
	sort.Strings(parentCandidates)
	parentCandidates = append([]string{PerennialBranchOption}, parentCandidates...)
	dialogData := enterParentModel{
		bubbleList: newBubbleList(parentCandidates, mainBranch.String()),
		branch:     branch.String(),
		mainBranch: mainBranch.String(),
	}
	dialogResult, err := tea.NewProgram(dialogData).Run()
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	result := dialogResult.(enterParentModel) //nolint:forcetypeassert
	selectedBranch := gitdomain.LocalBranchName(result.selectedEntry())
	return selectedBranch, result.aborted, nil
}

type enterParentModel struct {
	bubbleList
	branch     string
	mainBranch string
}

func (self enterParentModel) Init() tea.Cmd {
	return nil
} //nolint:ireturn

func (self enterParentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.bubbleList.handleKey(keyMsg); handled {
		return self, cmd
	}
	if keyMsg.Type == tea.KeyEnter {
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		return self, tea.Quit
	}
	return self, nil
}

func (self enterParentModel) View() string {
	s := strings.Builder{}
	s.WriteString("To sync branch \"" + self.branch + "\", Git Town needs to know its parent branch.\n")
	s.WriteString("Typically this is the main branch: " + self.mainBranch + "\n")
	s.WriteString("You can also select another feature or perennial branch.")
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
