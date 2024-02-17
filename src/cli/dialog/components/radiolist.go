package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v12/src/gohacks/slice"
)

// how many elements to display in the dialog
const WindowSize = 9

// RadioList lets the user select a new main branch for this repo.
func RadioList[S fmt.Stringer](entries []S, cursor int, title, help string, inputs TestInput) (selected S, aborted bool, err error) { //nolint:ireturn
	program := tea.NewProgram(radioListModel[S]{
		BubbleList: NewBubbleList(entries, cursor),
		help:       help,
		title:      title,
	})
	SendInputs(inputs, program)
	dialogResult, err := program.Run()
	if err != nil {
		return entries[0], false, err
	}
	result := dialogResult.(radioListModel[S]) //nolint:forcetypeassert
	return result.SelectedEntry(), result.Aborted(), nil
}

type radioListModel[S fmt.Stringer] struct {
	BubbleList[S]
	help  string // help text to display before the radio list
	title string // title to display before the help text
}

func (self radioListModel[S]) Init() tea.Cmd {
	return nil
}

func (self radioListModel[S]) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.BubbleList.HandleKey(keyMsg); handled {
		return self, cmd
	}
	if keyMsg.Type == tea.KeyEnter {
		self.Status = StatusDone
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self.Status = StatusDone
		return self, tea.Quit
	}
	return self, nil
}

func (self radioListModel[S]) View() string {
	if self.Status != StatusActive {
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
		s.WriteString(self.EntryNumberStr(i))
		if i == self.Cursor {
			s.WriteString(self.Colors.Selection.Styled("> " + branch.String()))
		} else {
			s.WriteString("  " + branch.String())
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
	s.WriteString(self.Colors.Help.Styled(" page up   "))
	// right
	s.WriteString(self.Colors.HelpKey.Styled("→"))
	s.WriteString(self.Colors.Help.Styled("/"))
	s.WriteString(self.Colors.HelpKey.Styled("d"))
	s.WriteString(self.Colors.Help.Styled(" page down   "))
	// numbers
	s.WriteString(self.Colors.HelpKey.Styled("0"))
	s.WriteString(self.Colors.Help.Styled("-"))
	s.WriteString(self.Colors.HelpKey.Styled("9"))
	s.WriteString(self.Colors.Help.Styled(" jump   "))
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
