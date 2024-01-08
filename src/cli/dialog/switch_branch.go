package dialog

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

type Model struct {
	Cursor        int
	Branches      []string
	CurrentBranch string
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("Grocery List")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func NewBuilder(lineage configdomain.Lineage) Builder {
	return Builder{
		Entries: ModalSelectEntries{},
		Lineage: lineage,
	}
}

// queryBranch lets the user select a new branch via a visual dialog.
// Indicates via `validSelection` whether the user made a valid selection.
func SwitchBranch(roots gitdomain.LocalBranchNames, selected gitdomain.LocalBranchName, lineage configdomain.Lineage) (selection gitdomain.LocalBranchName, validSelection bool, err error) {
	builder := NewBuilder(lineage)
	err = builder.CreateEntries(roots, selected)
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	choice, err := ModalSelect(builder.Entries, selected.String())
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	if choice == nil {
		return gitdomain.EmptyLocalBranchName(), false, nil
	}
	return gitdomain.NewLocalBranchName(*choice), true, nil
}

// Builder builds up the switch-branch dialog entries.
type Builder struct {
}

// AddEntryAndChildren adds the given branch and all its child branches to the given entries collection.
func (self *Builder) AddEntryAndChildren(branch gitdomain.LocalBranchName, indent int) error {
	self.Entries = append(self.Entries, ModalSelectEntry{
		Text:  strings.Repeat("  ", indent) + branch.String(),
		Value: branch.String(),
	})
	var err error
	for _, child := range self.Lineage.Children(branch) {
		err = self.AddEntryAndChildren(child, indent+1)
		if err != nil {
			return err
		}
	}
	return nil
}

// createEntries provides all the entries for the branch dialog.
func (self *Builder) CreateEntries(roots gitdomain.LocalBranchNames, selected gitdomain.LocalBranchName) error {
	var err error
	for _, root := range roots {
		err = self.AddEntryAndChildren(root, 0)
		if err != nil {
			return err
		}
	}
	if len(self.Entries) == 0 {
		self.Entries = append(self.Entries, ModalSelectEntry{
			Text:  string(selected),
			Value: string(selected),
		})
	}
	return nil
}
