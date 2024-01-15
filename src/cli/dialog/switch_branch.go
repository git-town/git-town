package dialog

import (
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

func SwitchBranch(localBranches gitdomain.LocalBranchNames, initialBranch gitdomain.LocalBranchName, lineage configdomain.Lineage) (gitdomain.LocalBranchName, bool, error) {
	entries := SwitchBranchEntries(localBranches, lineage)
	cursor := 0
	initialBranchName := initialBranch.String()
	for e, entry := range entries {
		if strings.TrimSpace(entry) == initialBranchName {
			cursor = e
		}
	}
	dialogData := SwitchModel{
		BubbleList:       newBubbleList(entries, cursor),
		InitialBranchPos: cursor,
	}
	dialogProcess := tea.NewProgram(dialogData, tea.WithOutput(os.Stderr))
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", false, err
	}
	result := dialogResult.(SwitchModel) //nolint:forcetypeassert
	selectedEntry := result.BubbleList.selectedEntry()
	selectedEntry = strings.TrimSpace(selectedEntry)
	selectedBranch := gitdomain.NewLocalBranchName(selectedEntry)
	return selectedBranch, result.Aborted, nil
}

type SwitchModel struct {
	BubbleList
	InitialBranchPos int // position of the currently checked out branch in the list
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
		case i == self.InitialBranchPos:
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

func SwitchBranchEntries(localBranches gitdomain.LocalBranchNames, lineage configdomain.Lineage) []string {
	entries := make([]string, 0, len(lineage))
	roots := lineage.Roots()
	// add all entries from the lineage
	for _, root := range roots {
		layoutBranches(&entries, root, "", lineage)
	}
	// remove entries for which no local branches exist

	// add missing local branches
	branchesInLineage := maps.Keys(lineage)
	for _, localBranch := range localBranches {
		if slices.Contains(roots, localBranch) {
			continue
		}
		if slices.Contains(branchesInLineage, localBranch) {
			continue
		}
		entries = append(entries, localBranch.String())
	}
	return entries
}

func layoutBranches(result *[]string, branch gitdomain.LocalBranchName, indentation string, lineage configdomain.Lineage) {
	*result = append(*result, indentation+branch.String())
	for _, child := range lineage.Children(branch) {
		layoutBranches(result, child, indentation+"  ", lineage)
	}
}
