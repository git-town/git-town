package dialog

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/muesli/termenv"
)

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(localBranches gitdomain.LocalBranchNames, oldPerennialBranches gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchNames, bool, error) {
	perennialCandidates := localBranches.Remove(mainBranch).AppendAllMissing(oldPerennialBranches)
	dialogData := perennialBranchesModel{
		bubbleList:    newBubbleList(perennialCandidates.Strings(), ""),
		selections:    slice.FindMany(perennialCandidates, oldPerennialBranches),
		selectedColor: termenv.String().Foreground(termenv.ANSIGreen),
	}
	dialogResult, err := tea.NewProgram(dialogData).Run()
	if err != nil {
		return gitdomain.LocalBranchNames{}, false, err
	}
	result := dialogResult.(perennialBranchesModel) //nolint:forcetypeassert
	selectedBranches := gitdomain.NewLocalBranchNames(result.selectedEntries()...)
	return selectedBranches, result.aborted, nil
}

type perennialBranchesModel struct {
	bubbleList
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
	if handled, cmd := self.bubbleList.handleKey(keyMsg); handled {
		return self, cmd
	}
	switch keyMsg.Type { //nolint:exhaustive
	case tea.KeySpace:
		self.toggleCurrentEntry()
		return self, nil
	case tea.KeyEnter:
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self.toggleCurrentEntry()
		return self, nil
	}
	return self, nil
}

func (self perennialBranchesModel) View() string {
	s := strings.Builder{}
	s.WriteString("Let's configure the perennial branches.\n")
	s.WriteString("These are long-lived branches, they are never shipped.\n")
	s.WriteString("Typically, perennial branches have names like \"development\", \"staging\", \"qa\", \"production\", etc.\n\n")
	for i, branch := range self.entries {
		selected := self.cursor == i
		checked := self.isRowChecked(i)
		switch {
		case selected && checked:
			s.WriteString(self.colors.selection.Styled("> [x] " + branch))
		case selected && !checked:
			s.WriteString(self.colors.selection.Styled("> [ ] " + branch))
		case !selected && checked:
			s.WriteString(self.selectedColor.Styled("  [x] " + branch))
		case !selected && !checked:
			s.WriteString("  [ ] " + branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(self.colors.helpKey.Styled("↑"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("k"))
	s.WriteString(self.colors.help.Styled(" up   "))
	// down
	s.WriteString(self.colors.helpKey.Styled("↓"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("j"))
	s.WriteString(self.colors.help.Styled(" down   "))
	// toggle
	s.WriteString(self.colors.helpKey.Styled("space"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("o"))
	s.WriteString(self.colors.help.Styled(" toggle   "))
	// accept
	s.WriteString(self.colors.helpKey.Styled("enter"))
	s.WriteString(self.colors.help.Styled(" accept   "))
	// abort
	s.WriteString(self.colors.helpKey.Styled("ctrl-c"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("q"))
	s.WriteString(self.colors.help.Styled(" abort"))
	return s.String()
}

func (self *perennialBranchesModel) disableCurrentEntry() {
	selectionIndex := slices.Index(self.selections, self.cursor)
	if selectionIndex != -1 {
		self.selections = slice.RemoveAt(self.selections, selectionIndex)
	}
}

func (self *perennialBranchesModel) enableCurrentEntry() {
	selectionIndex := slices.Index(self.selections, self.cursor)
	if selectionIndex == -1 {
		self.selections = append(self.selections, self.cursor)
	}
}

func (self *perennialBranchesModel) isRowChecked(row int) bool {
	return slices.Contains(self.selections, row)
}

func (self *perennialBranchesModel) selectedEntries() []string {
	result := []string{}
	for e, entry := range self.entries {
		if self.isRowChecked(e) {
			result = append(result, entry)
		}
	}
	return result
}

func (self *perennialBranchesModel) toggleCurrentEntry() {
	if self.isRowChecked(self.cursor) {
		self.disableCurrentEntry()
	} else {
		self.enableCurrentEntry()
	}
}
