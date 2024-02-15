package dialog

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/gohacks/slice"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/muesli/termenv"
)

const (
	perennialBranchesTitle = `Perennial branches`
	PerennialBranchesHelp  = `
Perennial branches are long-lived branches.
They are never shipped and have no ancestors.
Typically, perennial branches have names like
"development", "staging", "qa", "production", etc.

`
)

// PerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func PerennialBranches(localBranches gitdomain.LocalBranchNames, oldPerennialBranches gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName, inputs components.TestInput) (gitdomain.LocalBranchNames, bool, error) {
	perennialCandidates := localBranches.Remove(mainBranch).AppendAllMissing(oldPerennialBranches...)
	if len(perennialCandidates) == 0 {
		return gitdomain.LocalBranchNames{}, false, nil
	}
	program := tea.NewProgram(PerennialBranchesModel{
		BubbleList:    components.NewBubbleList(perennialCandidates, 0),
		Selections:    slice.FindMany(perennialCandidates, oldPerennialBranches),
		selectedColor: termenv.String().Foreground(termenv.ANSIGreen),
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
	components.BubbleList[gitdomain.LocalBranchName]
	Selections    []int
	selectedColor termenv.Style
}

// checkedEntries provides all checked list entries.
func (self *PerennialBranchesModel) CheckedEntries() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for e, entry := range self.Entries {
		if self.IsRowChecked(e) {
			result = append(result, entry)
		}
	}
	return result
}

// disableCurrentEntry unchecks the currently selected list entry.
func (self *PerennialBranchesModel) DisableCurrentEntry() {
	self.Selections = slice.Remove(self.Selections, self.Cursor)
}

// enableCurrentEntry checks the currently selected list entry.
func (self *PerennialBranchesModel) EnableCurrentEntry() {
	self.Selections = slice.AppendAllMissing(self.Selections, self.Cursor)
}

func (self PerennialBranchesModel) Init() tea.Cmd {
	return nil
}

// isRowChecked indicates whether the row with the given number is checked or not.
func (self *PerennialBranchesModel) IsRowChecked(row int) bool {
	return slices.Contains(self.Selections, row)
}

// isSelectedRowChecked indicates whether the currently selected list entry is checked or not.
func (self *PerennialBranchesModel) IsSelectedRowChecked() bool {
	return self.IsRowChecked(self.Cursor)
}

// toggleCurrentEntry unchecks the currently selected list entry if it is checked,
// and checks it if it is unchecked.
func (self *PerennialBranchesModel) ToggleCurrentEntry() {
	if self.IsRowChecked(self.Cursor) {
		self.DisableCurrentEntry()
	} else {
		self.EnableCurrentEntry()
	}
}

func (self PerennialBranchesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.BubbleList.HandleKey(keyMsg); handled {
		return self, cmd
	}
	switch keyMsg.Type { //nolint:exhaustive
	case tea.KeySpace:
		self.ToggleCurrentEntry()
		return self, nil
	case tea.KeyEnter:
		self.Status = components.StatusDone
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self.ToggleCurrentEntry()
		return self, nil
	}
	return self, nil
}

func (self PerennialBranchesModel) View() string {
	if self.Status != components.StatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteRune('\n')
	s.WriteString(self.Colors.Title.Styled(perennialBranchesTitle))
	s.WriteRune('\n')
	s.WriteString(PerennialBranchesHelp)
	window := slice.Window(slice.WindowArgs{
		ElementCount: len(self.Entries),
		CursorPos:    self.Cursor,
		WindowSize:   components.WindowSize,
	})
	for i := window.StartRow; i < window.EndRow; i++ {
		branch := self.Entries[i]
		selected := self.Cursor == i
		checked := self.IsRowChecked(i)
		s.WriteString(self.EntryNumberStr(i))
		switch {
		case selected && checked:
			s.WriteString(self.Colors.Selection.Styled("> [x] " + branch.String()))
		case selected && !checked:
			s.WriteString(self.Colors.Selection.Styled("> [ ] " + branch.String()))
		case !selected && checked:
			s.WriteString(self.selectedColor.Styled("  [x] " + branch.String()))
		case !selected && !checked:
			s.WriteString("  [ ] " + branch.String())
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
