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
	perennialBranchesTitle = `Perennial branches`
	PerennialBranchesHelp  = `
Perennial branches are long-lived branches.
They are never shipped and have no ancestors.
Typically, perennial branches have names like
"development", "staging", "qa", "production", etc.

See also the "perennial-regex" setting.

`
)

// TODO: extract a components.Checkboxes struct similar to components.RadioList that implements a generic checkbox list.

// PerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func PerennialBranches(localBranches gitdomain.LocalBranchNames, oldPerennialBranches gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName, inputs components.TestInput) (gitdomain.LocalBranchNames, bool, error) {
	perennialCandidates := localBranches.Remove(mainBranch).AppendAllMissing(oldPerennialBranches...)
	if len(perennialCandidates) == 0 {
		return gitdomain.LocalBranchNames{}, false, nil
	}
	program := tea.NewProgram(PerennialBranchesModel{
		List:          list.NewList(list.NewEntries(perennialCandidates...), 0),
		Selections:    slice.FindMany(perennialCandidates, oldPerennialBranches),
		selectedColor: colors.Green(),
	})
	components.SendInputs(inputs, program)
	dialogResult, err := program.Run()
	if err != nil {
		return gitdomain.LocalBranchNames{}, false, err
	}
	result := dialogResult.(PerennialBranchesModel) //nolint:forcetypeassert
	selectedBranches := result.CheckedEntries()
	selectionText := strings.Join(selectedBranches.Strings(), ", ")
	if selectionText == "" {
		selectionText = "(none)"
	}
	fmt.Printf(messages.PerennialBranches, components.FormattedSelection(selectionText, result.Aborted()))
	return selectedBranches, result.Aborted(), nil
}

type PerennialBranchesModel struct {
	list.List[gitdomain.LocalBranchName]
	Selections    []int
	selectedColor termenv.Style
}

// checkedEntries provides all checked list entries.
func (self PerennialBranchesModel) CheckedEntries() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for e, entry := range self.Entries {
		if self.IsRowChecked(e) {
			result = append(result, entry.Data)
		}
	}
	return result
}

// disableCurrentEntry unchecks the currently selected list entry.
func (self PerennialBranchesModel) DisableCurrentEntry() PerennialBranchesModel {
	self.Selections = slice.Remove(self.Selections, self.Cursor)
	return self
}

// enableCurrentEntry checks the currently selected list entry.
func (self PerennialBranchesModel) EnableCurrentEntry() PerennialBranchesModel {
	self.Selections = slice.AppendAllMissing(self.Selections, self.Cursor)
	return self
}

func (self PerennialBranchesModel) Init() tea.Cmd {
	return nil
}

// isRowChecked indicates whether the row with the given number is checked or not.
func (self PerennialBranchesModel) IsRowChecked(row int) bool {
	return slices.Contains(self.Selections, row)
}

// isSelectedRowChecked indicates whether the currently selected list entry is checked or not.
func (self PerennialBranchesModel) IsSelectedRowChecked() bool {
	return self.IsRowChecked(self.Cursor)
}

// toggleCurrentEntry unchecks the currently selected list entry if it is checked,
// and checks it if it is unchecked.
func (self PerennialBranchesModel) ToggleCurrentEntry() PerennialBranchesModel {
	if self.IsRowChecked(self.Cursor) {
		self = self.DisableCurrentEntry()
	} else {
		self = self.EnableCurrentEntry()
	}
	return self
}

func (self PerennialBranchesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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

func (self PerennialBranchesModel) View() string {
	if self.Status != list.StatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteRune('\n')
	s.WriteString(self.Colors.Title.Styled(perennialBranchesTitle))
	s.WriteRune('\n')
	s.WriteString(PerennialBranchesHelp)
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
