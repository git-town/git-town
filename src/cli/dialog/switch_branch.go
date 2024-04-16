package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/muesli/termenv"
	"golang.org/x/exp/maps"
)

func SwitchBranch(localBranches gitdomain.LocalBranchNames, initialBranch gitdomain.LocalBranchName, lineage configdomain.Lineage, allBranches gitdomain.BranchInfos, inputs components.TestInput) (gitdomain.LocalBranchName, bool, error) {
	entries := SwitchBranchEntries(localBranches, lineage, allBranches)
	cursor := SwitchBranchCursorPos(entries, initialBranch)
	dialogProgram := tea.NewProgram(SwitchModel{
		BubbleList:       components.NewBubbleList(entries, cursor),
		InitialBranchPos: cursor,
	})
	components.SendInputs(inputs, dialogProgram)
	dialogResult, err := dialogProgram.Run()
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
	window := slice.Window(slice.WindowArgs{
		CursorPos:    self.Cursor,
		ElementCount: len(self.Entries),
		WindowSize:   components.WindowSize,
	})
	for i := window.StartRow; i < window.EndRow; i++ {
		branch := self.Entries[i]
		isSelected := i == self.Cursor
		isInitial := i == self.InitialBranchPos
		isEnabled := branch.Enabled
		switch {
		case isSelected:
			color := self.Colors.Selection
			if !isEnabled {
				color = color.Faint()
			}
			s.WriteString(color.Styled("> " + branch.String()))
		case isInitial:
			color := self.Colors.Initial
			if !isEnabled {
				color = color.Faint()
			}
			s.WriteString(color.Styled("* " + branch.String()))
		default:
			color := termenv.String()
			if !isEnabled {
				color = color.Faint()
			}
			s.WriteString(color.Styled("  " + branch.String()))
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
	// left
	s.WriteString(self.Colors.HelpKey.Styled("←"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("u"))
	s.WriteString(self.Colors.Help.Styled(" 10 up   "))
	// right
	s.WriteString(self.Colors.HelpKey.Styled("→"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("d"))
	s.WriteString(self.Colors.Help.Styled(" 10 down   "))
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
func SwitchBranchEntries(localBranches gitdomain.LocalBranchNames, lineage configdomain.Lineage, allBranches gitdomain.BranchInfos) []SwitchBranchEntry {
	entries := make([]SwitchBranchEntry, 0, len(lineage))
	roots := lineage.Roots()
	// add all entries from the lineage
	for _, root := range roots {
		layoutBranches(&entries, root, "", lineage, allBranches)
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
func layoutBranches(result *[]SwitchBranchEntry, branch gitdomain.LocalBranchName, indentation string, lineage configdomain.Lineage, allBranches gitdomain.BranchInfos) {
	if allBranches.HasLocalBranch(branch) || allBranches.HasMatchingTrackingBranchFor(branch) {
		branchInfo := allBranches.FindByLocalName(branch)
		inThisWorktree := branchInfo.SyncStatus != gitdomain.SyncStatusOtherWorktree
		*result = append(*result, SwitchBranchEntry{Branch: branch, Indentation: indentation, Enabled: inThisWorktree})
	}
	for _, child := range lineage.Children(branch) {
		layoutBranches(result, child, indentation+"  ", lineage, allBranches)
	}
}

type SwitchBranchEntry struct {
	Branch      gitdomain.LocalBranchName
	Indentation string
	Enabled     bool
}

func (sbe SwitchBranchEntry) String() string {
	return sbe.Indentation + sbe.Branch.String()
}

func (sbe SwitchBranchEntry) IsEnabled() bool {
	return sbe.Enabled
}
