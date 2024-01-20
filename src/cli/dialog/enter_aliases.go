package dialog

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/muesli/termenv"
)

const enterAliasesHelp = `
You can create shorter aliases for frequently used Git Town commands.
For example, if the "git town sync" command is aliased,
you can call it as "git sync".

Please select which Git Town commands should be shortened.
If you are not sure, select all :)

`

// EnterAliases lets the select which Git Town commands should have shorter aliases.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterAliases(all, selected configdomain.AliasableCommands, dialogTestInput TestInput) (configdomain.AliasableCommands, bool, error) {
	dialogData := AliasesModel{
		BubbleList:    newBubbleList(all.Strings(), 0),
		Selections:    slice.FindMany(all, selected),
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
		return []configdomain.AliasableCommand{}, false, err
	}
	result := dialogResult.(AliasesModel) //nolint:forcetypeassert
	selectedCommands := configdomain.NewAliasableCommands(result.checkedEntries()...)
	aborted := result.Status == dialogStatusAborted
	var selectionText string
	switch len(selectedCommands) {
	case 0:
		selectionText = "(none)"
	case len(configdomain.AllAliasableCommands()):
		selectionText = "(all)"
	default:
		selectionText = strings.Join(result.checkedEntries(), ", ")
	}
	fmt.Printf("Aliased commands: %s\n", formattedSelection(selectionText, aborted))
	return selectedCommands, aborted, nil
}

type AliasesModel struct {
	BubbleList
	Selections    []int
	selectedColor termenv.Style
}

func (self AliasesModel) Init() tea.Cmd {
	return nil
}

func (self AliasesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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
	switch keyMsg.String() {
	case "a":
		self.SelectAll()
	case "n":
		self.selectNone()
	case "o":
		self.Status = dialogStatusDone
		self.toggleCurrentEntry()
		return self, nil
	}
	return self, nil
}

func (self AliasesModel) View() string {
	if self.Status != dialogStatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteString(enterAliasesHelp)
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
	// select all/none
	s.WriteString(self.Colors.helpKey.Styled("a"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("n"))
	s.WriteString(self.Colors.help.Styled(" select all/none   "))
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
func (self *AliasesModel) checkedEntries() []string {
	result := []string{}
	for e, entry := range self.Entries {
		if self.isRowChecked(e) {
			result = append(result, entry)
		}
	}
	return result
}

// disableCurrentEntry unchecks the currently selected list entry.
func (self *AliasesModel) disableCurrentEntry() {
	self.Selections = slice.Remove(self.Selections, self.Cursor)
}

// enableCurrentEntry checks the currently selected list entry.
func (self *AliasesModel) enableCurrentEntry() {
	self.Selections = slice.AppendAllMissing(self.Selections, self.Cursor)
}

// isRowChecked indicates whether the row with the given number is checked or not.
func (self *AliasesModel) isRowChecked(row int) bool {
	return slices.Contains(self.Selections, row)
}

// isSelectedRowChecked indicates whether the currently selected list entry is checked or not.
func (self *AliasesModel) isSelectedRowChecked() bool {
	return self.isRowChecked(self.Cursor)
}

// checks all entries in the list
func (self *AliasesModel) SelectAll() {
	count := len(self.Entries)
	self.Selections = make([]int, count)
	for i := 0; i < count; i++ {
		self.Selections[i] = i
	}
}

// checks all entries in the list
func (self *AliasesModel) selectNone() {
	self.Selections = []int{}
}

// toggleCurrentEntry unchecks the currently selected list entry if it is checked,
// and checks it if it is unchecked.
func (self *AliasesModel) toggleCurrentEntry() {
	if self.isRowChecked(self.Cursor) {
		self.disableCurrentEntry()
	} else {
		self.enableCurrentEntry()
	}
}
