package dialog

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v16/internal/cli/colors"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/slice"
	"github.com/git-town/git-town/v16/internal/messages"
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

func SwitchBranch(entries []SwitchBranchEntry, cursor int, uncommittedChanges bool, displayTypes configdomain.DisplayTypes, inputs components.TestInput) (gitdomain.LocalBranchName, bool, error) {
	dialogProgram := tea.NewProgram(SwitchModel{
		InitialBranchPos:   cursor,
		List:               list.NewList(newSwitchBranchListEntries(entries), cursor),
		UncommittedChanges: uncommittedChanges,
		DisplayBranchTypes: displayTypes,
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
	DisplayBranchTypes configdomain.DisplayTypes
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
		if self.DisplayBranchTypes {
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
