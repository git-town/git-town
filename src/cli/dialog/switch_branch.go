package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

func SwitchBranch(localBranches gitdomain.LocalBranchNames, initialBranch gitdomain.LocalBranchName, lineage configdomain.Lineage) (gitdomain.LocalBranchName, bool, error) {
	entries := SwitchBranchEntries(localBranches, lineage)
	cursor := SwitchBranchCursorPos(entries, initialBranch)
	dialogData := SwitchModel{
		BubbleList:       newBubbleList(entries, cursor),
		InitialBranchPos: cursor,
	}
	dialogProcess := tea.NewProgram(dialogData)
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", false, err
	}
	result := dialogResult.(SwitchModel) //nolint:forcetypeassert
	selectedEntry := result.BubbleList.selectedEntry()
	selectedEntry = strings.TrimSpace(selectedEntry)
	selectedBranch := gitdomain.NewLocalBranchName(selectedEntry)
	return selectedBranch, result.Status == dialogStatusAborted, nil
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
		self.Status = dialogStatusDone
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self.Status = dialogStatusDone
		return self, tea.Quit
	}
	return self, nil
}

func (self SwitchModel) View() string {
	if self.Status != dialogStatusActive {
		return ""
	}
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
	s.WriteString(self.Colors.helpKey.Styled("q"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("esc"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("ctrl-c"))
	s.WriteString(self.Colors.help.Styled(" abort"))
	return s.String()
}

// SwitchBranchCursorPos provides the initial cursor position for the "switch branch" dialog.
func SwitchBranchCursorPos(entries []string, initialBranch gitdomain.LocalBranchName) int {
	initialBranchName := initialBranch.String()
	for e, entry := range entries {
		if strings.TrimSpace(entry) == initialBranchName {
			return e
		}
	}
	return 0
}

// SwitchBranchEntries provides the entries for the "switch branch" dialog.
func SwitchBranchEntries(localBranches gitdomain.LocalBranchNames, lineage configdomain.Lineage) []string {
	entries := make([]string, 0, len(lineage))
	roots := lineage.Roots()
	// add all entries from the lineage
	for _, root := range roots {
		layoutBranches(&entries, root, "", lineage)
	}
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

// layoutBranches adds entries for the given branch and its children to the given entry list.
// The entries are indented according to their position in the given lineage.
func layoutBranches(result *[]string, branch gitdomain.LocalBranchName, indentation string, lineage configdomain.Lineage) {
	*result = append(*result, indentation+branch.String())
	for _, child := range lineage.Children(branch) {
		layoutBranches(result, child, indentation+"  ", lineage)
	}
}
