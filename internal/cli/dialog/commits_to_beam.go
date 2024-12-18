package dialog

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v17/internal/cli/colors"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks/slice"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/muesli/termenv"
)

const (
	commitsToBeamTitle = `Select the commits to beam into branch %s`
)

// PerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func CommitsToBeam(commits []gitdomain.Commit, targetBranch gitdomain.LocalBranchName, inputs components.TestInput) (gitdomain.Commits, bool, error) {
	entries := make([]commitsToBeamEntry, len(commits))
	for c, commit := range commits {
		entries[c] = commitsToBeamEntry(commit)
	}
	program := tea.NewProgram(commitsToBeamModel{
		List:          list.NewList(list.NewEntries(entries...), 0),
		Selections:    []int{},
		selectedColor: colors.Green(),
		targetBranch:  targetBranch,
	})
	components.SendInputs(inputs, program)
	dialogResult, err := program.Run()
	if err != nil {
		return gitdomain.Commits{}, false, err
	}
	result := dialogResult.(commitsToBeamModel) //nolint:forcetypeassert
	selectedBranches := result.CheckedEntries()
	fmt.Printf(messages.CommitsSelected, len(selectedBranches))
	return selectedBranches, result.Aborted(), nil
}

type commitsToBeamEntry gitdomain.Commit

func (entry commitsToBeamEntry) String() string {
	return entry.Message.String()
}

type commitsToBeamModel struct {
	list.List[commitsToBeamEntry]
	Selections    []int
	selectedColor termenv.Style
	targetBranch  gitdomain.LocalBranchName
}

// checkedEntries provides all checked list entries.
func (self commitsToBeamModel) CheckedEntries() []gitdomain.Commit {
	result := []gitdomain.Commit{}
	for e, entry := range self.Entries {
		if self.IsRowChecked(e) {
			result = append(result, gitdomain.Commit(entry.Data))
		}
	}
	return result
}

// disableCurrentEntry unchecks the currently selected list entry.
func (self commitsToBeamModel) DisableCurrentEntry() commitsToBeamModel {
	self.Selections = slice.Remove(self.Selections, self.Cursor)
	return self
}

// enableCurrentEntry checks the currently selected list entry.
func (self commitsToBeamModel) EnableCurrentEntry() commitsToBeamModel {
	self.Selections = slice.AppendAllMissing(self.Selections, self.Cursor)
	return self
}

func (model commitsToBeamModel) Init() tea.Cmd {
	return nil
}

// isRowChecked indicates whether the row with the given number is checked or not.
func (self commitsToBeamModel) IsRowChecked(row int) bool {
	return slices.Contains(self.Selections, row)
}

// isSelectedRowChecked indicates whether the currently selected list entry is checked or not.
func (self commitsToBeamModel) IsSelectedRowChecked() bool {
	return self.IsRowChecked(self.Cursor)
}

// toggleCurrentEntry unchecks the currently selected list entry if it is checked,
// and checks it if it is unchecked.
func (self commitsToBeamModel) ToggleCurrentEntry() commitsToBeamModel {
	if self.IsRowChecked(self.Cursor) {
		self = self.DisableCurrentEntry()
	} else {
		self = self.EnableCurrentEntry()
	}
	return self
}

func (self commitsToBeamModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.List.HandleKey(keyMsg); handled {
		return self, cmd
	}
	switch keyMsg.Type { //nolint:exhaustive
	case tea.KeySpace:
		self = self.ToggleCurrentEntry()
		return self, nil
	case tea.KeyEnter:
		self.Status = list.StatusDone
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self = self.ToggleCurrentEntry()
		return self, nil
	}
	return self, nil
}

func (self commitsToBeamModel) View() string {
	if self.Status != list.StatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteRune('\n')
	s.WriteString(self.Colors.Title.Styled(fmt.Sprintf(commitsToBeamTitle, self.selectedColor.Styled(self.targetBranch.String()))))
	s.WriteRune('\n')
	window := slice.Window(slice.WindowArgs{
		CursorPos:    self.Cursor,
		ElementCount: len(self.Entries),
		WindowSize:   components.WindowSize,
	})
	for i := window.StartRow; i < window.EndRow; i++ {
		branch := self.Entries[i]
		selected := self.Cursor == i
		checked := self.IsRowChecked(i)
		s.WriteString(self.EntryNumberStr(i))
		switch {
		case selected && checked:
			s.WriteString(self.Colors.Selection.Styled("> [x] " + branch.Text))
		case selected && !checked:
			s.WriteString(self.Colors.Selection.Styled("> [ ] " + branch.Text))
		case !selected && checked:
			s.WriteString(self.selectedColor.Styled("  [x] " + branch.Text))
		case !selected && !checked:
			s.WriteString("  [ ] " + branch.Text)
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
	// toggle
	s.WriteString(self.Colors.HelpKey.Styled("space"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("o"))
	s.WriteString(self.Colors.Help.Styled(" toggle   "))
	// numbers
	s.WriteString(self.Colors.HelpKey.Styled("0"))
	s.WriteString(self.Colors.Help.Styled("-"))
	s.WriteString(self.Colors.HelpKey.Styled("9"))
	s.WriteString(self.Colors.Help.Styled(" jump   "))
	// accept
	s.WriteString(self.Colors.HelpKey.Styled("enter"))
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
