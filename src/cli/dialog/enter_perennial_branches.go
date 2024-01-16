package dialog

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/muesli/termenv"
)

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(localBranches gitdomain.LocalBranchNames, oldPerennialBranches gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName, dialogTestInput TestInput) (gitdomain.LocalBranchNames, bool, error) {
	perennialCandidates := localBranches.Remove(mainBranch).AppendAllMissing(oldPerennialBranches...)
	if len(perennialCandidates) == 0 {
		return gitdomain.LocalBranchNames{}, false, nil
	}
	dialogData := perennialBranchesModel{
		BubbleList:    newBubbleList(perennialCandidates.Strings(), 0),
		selections:    slice.FindMany(perennialCandidates, oldPerennialBranches),
		selectedColor: termenv.String().Foreground(termenv.ANSIGreen),
	}
	program := tea.NewProgram(dialogData)
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
	result := dialogResult.(perennialBranchesModel) //nolint:forcetypeassert
	selectedBranches := gitdomain.NewLocalBranchNames(result.checkedEntries()...)
	aborted := result.Status == dialogStatusAborted
	fmt.Printf("Selected perennial branches: %s\n", formattedSelection(strings.Join(result.checkedEntries(), ", "), aborted))
	return selectedBranches, aborted, nil
}

type perennialBranchesModel struct {
	BubbleList
	selections    []int
	selectedColor termenv.Style
}

func (self perennialBranchesModel) Init() tea.Cmd {
	return nil
}

func (self perennialBranchesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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

func (self perennialBranchesModel) View() string {
	if self.Status != dialogStatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteString("Let's configure the perennial branches.\n")
	s.WriteString("These are long-lived branches without ancestors and are never shipped.\n")
	s.WriteString("Typically, perennial branches have names like \"development\", \"staging\", \"qa\", \"production\", etc.\n\n")
	for i, branch := range self.Entries {
		selected := self.Cursor == i
		checked := self.isRowChecked(i)
		s.WriteString(self.entryNumberStr(i))
		switch {
		case selected && checked:
			s.WriteString(self.Colors.selection.Styled("> [x] " + branch))
		case selected && !checked:
			s.WriteString(self.Colors.selection.Styled("> [ ] " + branch))
		case !selected && checked:
			s.WriteString(self.selectedColor.Styled("  [x] " + branch))
		case !selected && !checked:
			s.WriteString("  [ ] " + branch)
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
func (self *perennialBranchesModel) checkedEntries() []string {
	result := []string{}
	for e, entry := range self.Entries {
		if self.isRowChecked(e) {
			result = append(result, entry)
		}
	}
	return result
}

// disableCurrentEntry unchecks the currently selected list entry.
func (self *perennialBranchesModel) disableCurrentEntry() {
	self.selections = slice.Remove(self.selections, self.Cursor)
}

// enableCurrentEntry checks the currently selected list entry.
func (self *perennialBranchesModel) enableCurrentEntry() {
	self.selections = slice.AppendAllMissing(self.selections, self.Cursor)
}

// isRowChecked indicates whether the row with the given number is checked or not.
func (self *perennialBranchesModel) isRowChecked(row int) bool {
	return slices.Contains(self.selections, row)
}

// isSelectedRowChecked indicates whether the currently selected list entry is checked or not.
func (self *perennialBranchesModel) isSelectedRowChecked() bool {
	return self.isRowChecked(self.Cursor)
}

// toggleCurrentEntry unchecks the currently selected list entry if it is checked,
// and checks it if it is unchecked.
func (self *perennialBranchesModel) toggleCurrentEntry() {
	if self.isRowChecked(self.Cursor) {
		self.disableCurrentEntry()
	} else {
		self.enableCurrentEntry()
	}
}
