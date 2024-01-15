package dialog

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

func SwitchBranch(localBranches gitdomain.LocalBranchNames, initialBranch gitdomain.LocalBranchName, lineage configdomain.Lineage) (gitdomain.LocalBranchName, error) {
	entries := make([]string, 0, len(lineage))
	for _, root := range lineage.Roots() {
		layoutBranches(&entries, root, "", lineage)
	}
	dialogData := SwitchModel{
		BubbleList: newBubbleList(entries, initialBranch.String()),
		InitialPos: initialBranch.String(),
	}
	dialogProcess := tea.NewProgram(dialogData, tea.WithOutput(os.Stderr))
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", err
	}
	result := dialogResult.(SwitchModel) //nolint:forcetypeassert
	selectedEntry := result.BubbleList.selectedEntry()
	selectedEntry = strings.TrimSpace(selectedEntry)
	selectedBranch := gitdomain.NewLocalBranchName(selectedEntry)
	return selectedBranch, nil
}

func layoutBranches(result *[]string, branch gitdomain.LocalBranchName, indentation string, lineage configdomain.Lineage) {
	*result = append(*result, indentation+branch.String())
	for _, child := range lineage.Children(branch) {
		layoutBranches(result, child, indentation+"  ", lineage)
	}
}

type SwitchModel struct {
	BubbleList
	InitialBranch string // name of the currently checked out branch
}

func (self SwitchModel) Init() tea.Cmd {
	return nil
}

func (self SwitchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, code := self.BubbleList.handleKey(keyMsg); handled {
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

func (self SwitchModel) View() string {
	s := strings.Builder{}
	for i, branch := range self.Entries {
		switch {
		case i == self.Cursor:
			s.WriteString(self.Colors.selection.Styled("> " + branch))
		case branch == self.InitialBranch:
			s.WriteString(self.Colors.initial.Styled("* " + branch))
		default:
			s.WriteString("  " + branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(self.Colors.helpKey.Styled("↑"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("k"))
	s.WriteString(self.Colors.help.Styled(" up   "))
	// down
	s.WriteString(self.Colors.helpKey.Styled("↓"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("j"))
	s.WriteString(self.Colors.help.Styled(" down   "))
	// accept
	s.WriteString(self.Colors.helpKey.Styled("enter"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("o"))
	s.WriteString(self.Colors.help.Styled(" accept   "))
	// abort
	s.WriteString(self.Colors.helpKey.Styled("esc"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("q"))
	s.WriteString(self.Colors.help.Styled(" abort"))
	return s.String()
}
