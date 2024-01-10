package dialog

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"golang.org/x/term"
)

const PerennialBranchOption = "<none> (perennial branch)"

// EnterParent lets the user select the parent branch for the given branch.
func EnterParent(branch gitdomain.LocalBranchName, localBranches gitdomain.LocalBranchNames, lineage configdomain.Lineage, mainBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchName, bool, error) {
	parentCandidates := EnterParentEntries(branch, localBranches, lineage, mainBranch)
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return "", false, err
	}
	fmt.Println("term width:", termWidth)
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
	s.WriteString("\nPlease select the parent of branch \"" + self.branch + "\" or enter its number.\n")
	s.WriteString("Most of the time this is the main development branch (" + self.mainBranch + ").\n\n")
	for i, branch := range self.entries {
		if i == self.cursor {
			// TODO: display single or double-digit numbers for branches,
			// and also allow the user to enter the branch number.
			// If the number is double digits, the user must press two numbers.
			// This provides accessibility out of the box.
			s.WriteString(self.colors.selection.Styled(strconv.FormatInt(int64(i), 10) + " > " + branch))
		} else {
			s.WriteString(strconv.FormatInt(int64(i), 10) + "   " + branch)
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
	s.WriteString(self.colors.helpKey.Styled("ctrl-c"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("q"))
	s.WriteString(self.colors.help.Styled(" abort"))
	return s.String()
}

func EnterParentEntries(branch gitdomain.LocalBranchName, localBranches gitdomain.LocalBranchNames, lineage configdomain.Lineage, mainBranch gitdomain.LocalBranchName) []string {
	parentCandidateBranches := localBranches.Remove(branch).Remove(lineage.Children(branch)...)
	parentCandidateBranches.Sort()
	parentCandidates := parentCandidateBranches.Hoist(mainBranch).Strings()
	return append([]string{PerennialBranchOption}, parentCandidates...)
}
