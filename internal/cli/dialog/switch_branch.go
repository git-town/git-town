package dialog

import (
	"regexp"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/regexes"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/muesli/termenv"
)

type SwitchBranchEntry struct {
	Branch        gitdomain.LocalBranchName
	Indentation   string
	OtherWorktree bool
	Type          configdomain.BranchType
}

func (sbe SwitchBranchEntry) String() string {
	return sbe.Indentation + sbe.Branch.String()
}

type SwitchBranchEntries []SwitchBranchEntry

func (sbes SwitchBranchEntries) ContainsBranch(branch gitdomain.LocalBranchName) bool {
	for _, entry := range sbes {
		if entry.Branch == branch {
			return true
		}
	}
	return false
}

func (sbes SwitchBranchEntries) IndexOf(branch gitdomain.LocalBranchName) int {
	for e, entry := range sbes {
		if entry.Branch == branch {
			return e
		}
	}
	return 0
}

type SwitchModel struct {
	list.List[SwitchBranchEntry]
	DisplayBranchTypes configdomain.DisplayTypes
	InitialBranchPos   Option[int] // position of the currently checked out branch in the list
	UncommittedChanges bool        // whether the workspace has uncommitted changes
}

func (self SwitchModel) Init() tea.Cmd {
	return nil
}

func (self SwitchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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
	s := strings.Builder{}
	if self.Status != list.StatusActive {
		return ""
	}
	if self.UncommittedChanges {
		s.WriteString("\n")
		s.WriteString(colors.BoldCyan().Styled(messages.SwitchUncommittedChanges))
		s.WriteString("\n\n")
	}
	window := slice.Window(slice.WindowArgs{
		CursorPos:    self.Cursor,
		ElementCount: len(self.Entries),
		WindowSize:   dialogcomponents.WindowSize,
	})
	for i := window.StartRow; i < window.EndRow; i++ {
		entry := self.Entries[i]
		isSelected := i == self.Cursor
		initialBranchPos, hasInitialBranchPos := self.InitialBranchPos.Get()
		isInitial := hasInitialBranchPos && i == initialBranchPos
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
		if self.DisplayBranchTypes.IsTrue() && ShouldDisplayBranchType(entry.Data.Type) {
			s.WriteString("  ")
			s.WriteString(colors.Faint().Styled("(" + entry.Data.Type.String() + ")"))
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

// NewSwitchBranchEntries provides the entries for the "switch branch" components.
func NewSwitchBranchEntries(args NewSwitchBranchEntriesArgs) SwitchBranchEntries {
	entries := make(SwitchBranchEntries, 0, args.Lineage.Len())
	roots := args.Lineage.Roots()
	// add all entries from the lineage
	for _, root := range roots {
		layoutBranches(layoutBranchesArgs{
			branch:            root,
			branchInfos:       args.BranchInfos,
			branchTypes:       args.BranchTypes,
			branchesAndTypes:  args.BranchesAndTypes,
			excludeBranches:   args.ExcludeBranches,
			indentation:       "",
			lineage:           args.Lineage,
			regexes:           args.Regexes,
			result:            &entries,
			showAllBranches:   args.ShowAllBranches,
			unknownBranchType: args.UnknownBranchType,
		})
	}
	// add branches not in the lineage
	branchesInLineage := args.Lineage.BranchesWithParents()
	for _, branchInfo := range args.BranchInfos {
		localBranch := branchInfo.LocalBranchName()
		if slices.Contains(roots, localBranch) {
			continue
		}
		if slices.Contains(branchesInLineage, localBranch) {
			continue
		}
		if entries.ContainsBranch(localBranch) {
			continue
		}
		layoutBranches(layoutBranchesArgs{
			branch:            localBranch,
			branchInfos:       args.BranchInfos,
			branchTypes:       args.BranchTypes,
			branchesAndTypes:  args.BranchesAndTypes,
			excludeBranches:   args.ExcludeBranches,
			indentation:       "",
			lineage:           args.Lineage,
			regexes:           args.Regexes,
			result:            &entries,
			showAllBranches:   args.ShowAllBranches,
			unknownBranchType: args.UnknownBranchType,
		})
	}
	return entries
}

type NewSwitchBranchEntriesArgs struct {
	BranchInfos       gitdomain.BranchInfos
	BranchTypes       []configdomain.BranchType
	BranchesAndTypes  configdomain.BranchesAndTypes
	ExcludeBranches   gitdomain.LocalBranchNames
	Lineage           configdomain.Lineage
	Regexes           []*regexp.Regexp
	ShowAllBranches   configdomain.AllBranches
	UnknownBranchType configdomain.UnknownBranchType
}

func ShouldDisplayBranchType(branchType configdomain.BranchType) bool {
	switch branchType {
	case
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypeFeatureBranch:
		return false
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePerennialBranch,
		configdomain.BranchTypePrototypeBranch:
		return true
	}
	panic("unhandled branch type:" + branchType.String())
}

func SwitchBranch(args SwitchBranchArgs) (gitdomain.LocalBranchName, dialogdomain.Exit, error) {
	initialBranchPos := None[int]()
	if currentBranch, has := args.CurrentBranch.Get(); has {
		initialBranchPos = Some(args.Entries.IndexOf(currentBranch))
	}
	dialogProgram := tea.NewProgram(SwitchModel{
		DisplayBranchTypes: args.DisplayBranchTypes,
		InitialBranchPos:   initialBranchPos,
		List:               list.NewList(newSwitchBranchListEntries(args.Entries), args.Cursor),
		UncommittedChanges: args.UncommittedChanges,
	})
	dialogcomponents.SendInputs(args.InputName, args.Inputs.Next(), dialogProgram)
	dialogResult, err := dialogProgram.Run()
	result := dialogResult.(SwitchModel)
	selectedData := result.List.SelectedData()
	return selectedData.Branch, result.Aborted(), err
}

type SwitchBranchArgs struct {
	CurrentBranch      Option[gitdomain.LocalBranchName]
	Cursor             int
	DisplayBranchTypes configdomain.DisplayTypes
	Entries            SwitchBranchEntries
	InputName          string
	Inputs             dialogcomponents.Inputs
	UncommittedChanges bool
}

// layoutBranches adds entries for the given branch and its children to the given entry list.
// The entries are indented according to their position in the given lineage.
func layoutBranches(args layoutBranchesArgs) {
	if args.excludeBranches.Contains(args.branch) {
		return
	}
	if args.branchInfos.HasLocalBranch(args.branch) || args.showAllBranches.Enabled() {
		branchInfo, hasBranchInfo := args.branchInfos.FindByLocalName(args.branch).Get()
		otherWorktree := hasBranchInfo && branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
		branchType, hasBranchType := args.branchesAndTypes[args.branch]
		if !hasBranchType && len(args.branchTypes) > 0 {
			branchType = args.unknownBranchType.BranchType()
		}
		hasCorrectBranchType := len(args.branchTypes) == 0 || slices.Contains(args.branchTypes, branchType)
		matchesRegex := args.regexes.Matches(args.branch.String())
		if hasCorrectBranchType && matchesRegex {
			*args.result = append(*args.result, SwitchBranchEntry{
				Branch:        args.branch,
				Indentation:   args.indentation,
				OtherWorktree: otherWorktree,
				Type:          branchType,
			})
		}
	}
	for _, child := range args.lineage.Children(args.branch) {
		layoutBranches(layoutBranchesArgs{
			branch:            child,
			branchInfos:       args.branchInfos,
			branchTypes:       args.branchTypes,
			branchesAndTypes:  args.branchesAndTypes,
			excludeBranches:   args.excludeBranches,
			indentation:       args.indentation + "  ",
			lineage:           args.lineage,
			regexes:           args.regexes,
			result:            args.result,
			showAllBranches:   args.showAllBranches,
			unknownBranchType: args.unknownBranchType,
		})
	}
}

type layoutBranchesArgs struct {
	branch            gitdomain.LocalBranchName
	branchInfos       gitdomain.BranchInfos
	branchTypes       []configdomain.BranchType
	branchesAndTypes  configdomain.BranchesAndTypes
	excludeBranches   gitdomain.LocalBranchNames
	indentation       string
	lineage           configdomain.Lineage
	regexes           regexes.Regexes
	result            *SwitchBranchEntries
	showAllBranches   configdomain.AllBranches
	unknownBranchType configdomain.UnknownBranchType
}

func newSwitchBranchListEntries(switchBranchEntries []SwitchBranchEntry) list.Entries[SwitchBranchEntry] {
	result := make(list.Entries[SwitchBranchEntry], len(switchBranchEntries))
	for e, entry := range switchBranchEntries {
		result[e] = list.Entry[SwitchBranchEntry]{
			Data:     entry,
			Disabled: entry.OtherWorktree,
			Text:     entry.String(),
		}
	}
	return result
}
