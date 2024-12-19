package components

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/gohacks/slice"
)

// CheckList lets the user select zero to many of the given entries.
func CheckList[S comparable](entries list.Entries[S], title, help string, inputs TestInput) (selected []S, aborted bool, err error) { //nolint:ireturn
	program := tea.NewProgram(CheckListModel[S]{
		List:  list.NewList(entries, 0),
		help:  help,
		title: title,
	})
	SendInputs(inputs, program)
	dialogResult, err := program.Run()
	if err != nil {
		return []S{}, false, err
	}
	result := dialogResult.(CheckListModel[S]) //nolint:forcetypeassert
	return result.CheckedEntries(), result.Aborted(), nil
}

type CheckListModel[S comparable] struct {
	list.List[S]
	Selections []int
	help       string // help text to display before the checklist
	title      string // title to display before the help text
}

// checkedEntries provides all checked list entries.
func (self CheckListModel[S]) CheckedEntries() []S {
	result := []S{}
	for e, entry := range self.Entries {
		if self.IsRowChecked(e) {
			result = append(result, entry.Data)
		}
	}
	return result
}

// disableCurrentEntry unchecks the currently selected list entry.
func (self CheckListModel[S]) DisableCurrentEntry() CheckListModel[S] {
	self.Selections = slice.Remove(self.Selections, self.Cursor)
	return self
}

// enableCurrentEntry checks the currently selected list entry.
func (self CheckListModel[S]) EnableCurrentEntry() CheckListModel[S] {
	self.Selections = slice.AppendAllMissing(self.Selections, self.Cursor)
	return self
}

func (self CheckListModel[S]) Init() tea.Cmd {
	return nil
}

// isRowChecked indicates whether the row with the given number is checked or not.
func (self CheckListModel[S]) IsRowChecked(row int) bool {
	return slices.Contains(self.Selections, row)
}

// isSelectedRowChecked indicates whether the currently selected list entry is checked or not.
func (self CheckListModel[S]) IsSelectedRowChecked() bool {
	return self.IsRowChecked(self.Cursor)
}

// toggleCurrentEntry unchecks the currently selected list entry if it is checked,
// and checks it if it is unchecked.
func (self CheckListModel[S]) ToggleCurrentEntry() CheckListModel[S] {
	if self.IsRowChecked(self.Cursor) {
		self = self.DisableCurrentEntry()
	} else {
		self = self.EnableCurrentEntry()
	}
	return self
}

func (self CheckListModel[S]) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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

func (self CheckListModel[S]) View() string {
	if self.Status != list.StatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteRune('\n')
	s.WriteString(self.Colors.Title.Styled(self.title))
	s.WriteRune('\n')
	s.WriteString(self.help)
	window := slice.Window(slice.WindowArgs{
		CursorPos:    self.Cursor,
		ElementCount: len(self.Entries),
		WindowSize:   WindowSize,
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
			s.WriteString(self.Colors.Initial.Styled("  [x] " + branch.Text))
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
