package dialogscreens

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/muesli/termenv"
)

const enterPerennialBranchesHelp = `
	Perennial branches are long-lived branches.
	They are never shipped and don't have ancestors.
	Typically, perennial branches have names like
	"development", "staging", "qa", "production", etc.

`

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(localBranches gitdomain.LocalBranchNames, oldPerennialBranches gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName, dialogTestInput TestInput) (gitdomain.LocalBranchNames, bool, error) {
	perennialCandidates := localBranches.Remove(mainBranch).AppendAllMissing(oldPerennialBranches...)
	if len(perennialCandidates) == 0 {
		return gitdomain.LocalBranchNames{}, false, nil
	}
	program := tea.NewProgram(PerennialBranchesModel{
		BubbleList:    newBubbleList(perennialCandidates, 0),
		selections:    slice.FindMany(perennialCandidates, oldPerennialBranches),
		selectedColor: termenv.String().Foreground(termenv.ANSIGreen),
	})
	if len(dialogTestInput) > 0 {
		go func() {
			for _, input := range dialogTestInput {
				program.Send(input)
			}
		}()
	}
	dialogResult, err := program.Run()
	if err != nil {
		return gitdomain.LocalBranchNames{}, false, err
	}
	result := dialogResult.(PerennialBranchesModel) //nolint:forcetypeassert
	selectedBranches := result.checkedEntries()
	selectionText := strings.Join(selectedBranches.Strings(), ", ")
	if selectionText == "" {
		selectionText = "(none)"
	}
	fmt.Printf("Perennial branches: %s\n", formattedSelection(selectionText, result.aborted()))
	return selectedBranches, result.aborted(), nil
}

type PerennialBranchesModel struct {
	BubbleList[gitdomain.LocalBranchName]
	selections    []int
	selectedColor termenv.Style
}

func (self PerennialBranchesModel) Init() tea.Cmd {
	return nil
}

func (self PerennialBranchesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.BubbleList.handleKey(keyMsg); handled {
		return self, cmd
	}
	switch keyMsg.Type { //nolint:exhaustive
	case tea.KeySpace:
		self.toggleCurrentEntry()
		return self, nil
	case tea.KeyEnter:
		self.Status = dialogStatusDone
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self.Status = dialogStatusDone
		self.toggleCurrentEntry()
		return self, nil
	}
	return self, nil
}

func (self PerennialBranchesModel) View() string {
	if self.Status != dialogStatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteString(enterPerennialBranchesHelp)
	for i, branch := range self.Entries {
		selected := self.Cursor == i
		checked := self.isRowChecked(i)
		s.WriteString(self.entryNumberStr(i))
		switch {
		case selected && checked:
			s.WriteString(self.Colors.selection.Styled("> [x] " + branch.String()))
		case selected && !checked:
			s.WriteString(self.Colors.selection.Styled("> [ ] " + branch.String()))
		case !selected && checked:
			s.WriteString(self.selectedColor.Styled("  [x] " + branch.String()))
		case !selected && !checked:
			s.WriteString("  [ ] " + branch.String())
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
	// toggle
	s.WriteString(self.Colors.helpKey.Styled("space"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("o"))
	s.WriteString(self.Colors.help.Styled(" toggle   "))
	// numbers
	s.WriteString(self.Colors.helpKey.Styled("0"))
	s.WriteString(self.Colors.help.Styled("-"))
	s.WriteString(self.Colors.helpKey.Styled("9"))
	s.WriteString(self.Colors.help.Styled(" jump   "))
	// accept
	s.WriteString(self.Colors.helpKey.Styled("enter"))
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

// checkedEntries provides all checked list entries.
func (self *PerennialBranchesModel) checkedEntries() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for e, entry := range self.Entries {
		if self.isRowChecked(e) {
			result = append(result, entry)
		}
	}
	return result
}

// disableCurrentEntry unchecks the currently selected list entry.
func (self *PerennialBranchesModel) disableCurrentEntry() {
	self.selections = slice.Remove(self.selections, self.Cursor)
}

// enableCurrentEntry checks the currently selected list entry.
func (self *PerennialBranchesModel) enableCurrentEntry() {
	self.selections = slice.AppendAllMissing(self.selections, self.Cursor)
}

// isRowChecked indicates whether the row with the given number is checked or not.
func (self *PerennialBranchesModel) isRowChecked(row int) bool {
	return slices.Contains(self.selections, row)
}

// isSelectedRowChecked indicates whether the currently selected list entry is checked or not.
func (self *PerennialBranchesModel) isSelectedRowChecked() bool {
	return self.isRowChecked(self.Cursor)
}

// toggleCurrentEntry unchecks the currently selected list entry if it is checked,
// and checks it if it is unchecked.
func (self *PerennialBranchesModel) toggleCurrentEntry() {
	if self.isRowChecked(self.Cursor) {
		self.disableCurrentEntry()
	} else {
		self.enableCurrentEntry()
	}
}
