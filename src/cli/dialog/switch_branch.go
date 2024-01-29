package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

func SwitchBranch(localBranches gitdomain.LocalBranchNames, initialBranch gitdomain.LocalBranchName, lineage configdomain.Lineage) (gitdomain.LocalBranchName, bool, error) {
	entries := SwitchBranchEntries(localBranches, lineage)
	cursor := SwitchBranchCursorPos(entries, initialBranch)
	dialogProcess := tea.NewProgram(SwitchModel{
		BubbleList:       components.NewBubbleList(entries, cursor),
		InitialBranchPos: cursor,
	})
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", false, err
	}
	result := dialogResult.(SwitchModel) //nolint:forcetypeassert
	selectedEntry := result.BubbleList.SelectedEntry()
	return selectedEntry.Branch, result.Aborted(), nil
}

type SwitchModel struct {
	components.BubbleList[SwitchBranchEntry]
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
	if handled, code := self.BubbleList.HandleKey(keyMsg); handled {
		return self, code
	}
	if keyMsg.Type == tea.KeyEnter {
		self.Status = components.StatusDone
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self.Status = components.StatusDone
		return self, tea.Quit
	}
	return self, nil
}

func (self SwitchModel) View() string {
	if self.Status != components.StatusActive {
		return ""
	}
	s := strings.Builder{}
	for i, branch := range self.Entries {
		switch {
		case i == self.Cursor:
			s.WriteString(self.Colors.Selection.Styled("> " + branch.String()))
		case i == self.InitialBranchPos:
			s.WriteString(self.Colors.Initial.Styled("* " + branch.String()))
		default:
			s.WriteString("  " + branch.String())
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(self.Colors.HelpKey.Styled("↑"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("k"))
	s.WriteString(self.Colors.Help.Styled(" up   "))
	// down
	s.WriteString(self.Colors.HelpKey.Styled("↓"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("j"))
	s.WriteString(self.Colors.Help.Styled(" down   "))
	// accept
	s.WriteString(self.Colors.HelpKey.Styled("enter"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("o"))
	s.WriteString(self.Colors.Help.Styled(" accept   "))
	// abort
	s.WriteString(self.Colors.HelpKey.Styled("q"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("esc"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("ctrl-c"))
	s.WriteString(self.Colors.Help.Styled(" abort"))
	return s.String()
}

// SwitchBranchCursorPos provides the initial cursor position for the "switch branch" components.
func SwitchBranchCursorPos(entries []SwitchBranchEntry, initialBranch gitdomain.LocalBranchName) int {
	for e, entry := range entries {
		if entry.Branch == initialBranch {
			return e
		}
	}
	return 0
}

// SwitchBranchEntries provides the entries for the "switch branch" components.
func SwitchBranchEntries(localBranches gitdomain.LocalBranchNames, lineage configdomain.Lineage) []SwitchBranchEntry {
	entries := make([]SwitchBranchEntry, 0, len(lineage))
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
		entries = append(entries, SwitchBranchEntry{Branch: localBranch, Indentation: ""})
	}
	return entries
}

// layoutBranches adds entries for the given branch and its children to the given entry list.
// The entries are indented according to their position in the given lineage.
func layoutBranches(result *[]SwitchBranchEntry, branch gitdomain.LocalBranchName, indentation string, lineage configdomain.Lineage) {
	*result = append(*result, SwitchBranchEntry{Branch: branch, Indentation: indentation})
	for _, child := range lineage.Children(branch) {
		layoutBranches(result, child, indentation+"  ", lineage)
	}
}

type SwitchBranchEntry struct {
	Branch      gitdomain.LocalBranchName
	Indentation string
}

func (sbe SwitchBranchEntry) String() string {
	return sbe.Indentation + sbe.Branch.String()
}
