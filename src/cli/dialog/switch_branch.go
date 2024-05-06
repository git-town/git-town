package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/muesli/termenv"
	"golang.org/x/exp/maps"
)

type SwitchBranchEntry struct {
	Branch        gitdomain.LocalBranchName
	Indentation   string
	OtherWorktree bool
}

func (sbe SwitchBranchEntry) String() string {
	return sbe.Indentation + sbe.Branch.String()
}

func SwitchBranch(localBranches gitdomain.LocalBranchNames, initialBranch gitdomain.LocalBranchName, lineage configdomain.Lineage, allBranches gitdomain.BranchInfos, uncommittedChanges bool, inputs components.TestInput) (gitdomain.LocalBranchName, bool, error) {
	entries := SwitchBranchEntries(localBranches, lineage, allBranches)
	cursor := SwitchBranchCursorPos(entries, initialBranch)
	dialogProgram := tea.NewProgram(SwitchModel{
		InitialBranchPos:   cursor,
		List:               list.NewList(newSwitchBranchListEntries(entries), cursor),
		UncommittedChanges: uncommittedChanges,
	})
	components.SendInputs(inputs, dialogProgram)
	dialogResult, err := dialogProgram.Run()
	if err != nil {
		return "", false, err
	}
	result := dialogResult.(SwitchModel) //nolint:forcetypeassert
	selectedData := result.List.SelectedData()
	return selectedData.Branch, result.Aborted(), nil
}

type SwitchModel struct {
	list.List[SwitchBranchEntry]
	InitialBranchPos   int  // position of the currently checked out branch in the list
	UncommittedChanges bool // whether the workspace has uncommitted changes
}

func (self SwitchModel) Init() tea.Cmd {
	return nil
}

func (self SwitchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint: ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, code := self.List.HandleKey(keyMsg); handled {
		return self, code
	}
	if keyMsg.Type == tea.KeyEnter {
		self.Status = list.StatusDone
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self.Status = list.StatusDone
		return self, tea.Quit
	}
	return self, nil
}

func (self SwitchModel) View() string {
	if self.Status != list.StatusActive {
		return ""
	}
	s := strings.Builder{}
	if self.UncommittedChanges {
		s.WriteString("\n")
		s.WriteString(colors.BoldCyan().Styled(messages.SwitchUncommittedChanges))
		s.WriteString("\n")
	}
	window := slice.Window(slice.WindowArgs{
		CursorPos:    self.Cursor,
		ElementCount: len(self.Entries),
		WindowSize:   components.WindowSize,
	})
	for i := window.StartRow; i < window.EndRow; i++ {
		entry := self.Entries[i]
		isSelected := i == self.Cursor
		isInitial := i == self.InitialBranchPos
		switch {
		case isSelected:
			color := self.Colors.Selection
			if entry.Data.OtherWorktree {
				color = color.Faint()
			}
			s.WriteString(color.Styled("> " + entry.Text))
		case isInitial:
			color := self.Colors.Initial
			if entry.Data.OtherWorktree {
				color = color.Faint()
			}
			s.WriteString(color.Styled("* " + entry.Text))
		case entry.Data.OtherWorktree:
			s.WriteString(colors.Faint().Styled("+ " + entry.Text))
		default:
			color := termenv.String()
			if entry.Data.OtherWorktree {
				color = color.Faint()
			}
			s.WriteString(color.Styled("  " + entry.Text))
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
		var otherWorktree bool
		if branchInfo, hasBranchInfo := allBranches.FindByLocalName(localBranch).Get(); hasBranchInfo {
			otherWorktree = branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
		} else {
			otherWorktree = false
		}
		entries = append(entries, SwitchBranchEntry{Branch: localBranch, Indentation: "", OtherWorktree: otherWorktree})
	}
	return entries
}

// layoutBranches adds entries for the given branch and its children to the given entry list.
// The entries are indented according to their position in the given lineage.
func layoutBranches(result *[]SwitchBranchEntry, branch gitdomain.LocalBranchName, indentation string, lineage configdomain.Lineage, allBranches gitdomain.BranchInfos) {
	if allBranches.HasLocalBranch(branch) || allBranches.HasMatchingTrackingBranchFor(branch) {
		var otherWorktree bool
		if branchInfo, hasBranchInfo := allBranches.FindByLocalName(branch).Get(); hasBranchInfo {
			otherWorktree = branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
		} else {
			otherWorktree = false
		}
		*result = append(*result, SwitchBranchEntry{Branch: branch, Indentation: indentation, OtherWorktree: otherWorktree})
	}
	for _, child := range lineage.Children(branch) {
		layoutBranches(result, child, indentation+"  ", lineage, allBranches)
	}
}

func newSwitchBranchListEntries(switchBranchEntries []SwitchBranchEntry) list.Entries[SwitchBranchEntry] {
	result := make(list.Entries[SwitchBranchEntry], len(switchBranchEntries))
	for e, entry := range switchBranchEntries {
		result[e] = list.Entry[SwitchBranchEntry]{
			Data:    entry,
			Enabled: !entry.OtherWorktree,
			Text:    entry.String(),
		}
	}
	return result
}
