package dialog

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// EnterMainBranch lets the user select a new main branch for this repo.
func radioList(args radioListArgs) (selected string, aborted bool, err error) {
	model := radioListModel{
		bubbleList: newBubbleList(args.entries, args.defaultEntry),
		help:       args.help,
	}
	program := tea.NewProgram(model)
	inputText, hasInput := os.LookupEnv(TestInputKey)
	if hasInput {
		inputs := ParseTestInput(inputText)
		go func() {
			for _, input := range inputs {
				program.Send(input)
			}
		}()
	}
	dialogResult, err := program.Run()
	if err != nil {
		return "", false, err
	}
	result := dialogResult.(radioListModel) //nolint:forcetypeassert
	return result.selectedEntry(), result.aborted, nil
}

type radioListArgs struct {
	entries      []string
	defaultEntry string
	help         string
}

type radioListModel struct {
	bubbleList
	help string // help text to display before the radio list
}

func (self radioListModel) Init() tea.Cmd {
	return nil
}

func (self radioListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	fmt.Printf("RECEIVED MSG %#v\n", msg)
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.bubbleList.handleKey(keyMsg); handled {
		return self, cmd
	}
	if keyMsg.Type == tea.KeyEnter {
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		return self, tea.Quit
	}
	return self, nil
}

func (self radioListModel) View() string {
	s := strings.Builder{}
	s.WriteString(self.help)
	for i, branch := range self.entries {
		s.WriteString(self.entryNumberStr(i))
		if i == self.cursor {
			s.WriteString(self.colors.selection.Styled("> " + branch))
		} else {
			s.WriteString("  " + branch)
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
	// numbers
	s.WriteString(self.colors.helpKey.Styled("0"))
	s.WriteString(self.colors.help.Styled("-"))
	s.WriteString(self.colors.helpKey.Styled("9"))
	s.WriteString(self.colors.help.Styled(" jump   "))
	// accept
	s.WriteString(self.colors.helpKey.Styled("enter"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("o"))
	s.WriteString(self.colors.help.Styled(" accept   "))
	// abort
	s.WriteString(self.colors.helpKey.Styled("esc"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("q"))
	s.WriteString(self.colors.help.Styled(" abort"))
	return s.String()
}
