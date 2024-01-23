package dialog

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/muesli/termenv"
)

const enterAliasesHelp = `
You can create shorter aliases for frequently used Git Town commands.
For example, if the "git town sync" command is aliased,
you can call it as "git sync".

Please select which Git Town commands should be shortened.
If you are not sure, select all :)

`

// Aliases lets the user select which Git Town commands should have shorter aliases.
// This includes asking the user and updating the respective settings based on the user selection.
func Aliases(allAliasableCommands configdomain.AliasableCommands, existingAliases configdomain.Aliases, dialogTestInput TestInput) (configdomain.Aliases, bool, error) {
	program := tea.NewProgram(AliasesModel{
		AllAliasableCommands: allAliasableCommands,
		BubbleList:           newBubbleList(allAliasableCommands.Strings(), 0),
		CurrentSelections:    NewAliasSelections(allAliasableCommands, existingAliases),
		OriginalAliases:      existingAliases,
		selectedColor:        termenv.String().Foreground(termenv.ANSIGreen),
	})
	if len(dialogTestInput) > 0 {
		go func() {
			for _, input := range dialogTestInput {
				program.Send(input)
			}
		}()
	}
	dialogResult, err := program.Run()
	result := dialogResult.(AliasesModel) //nolint:forcetypeassert
	if err != nil || result.aborted() {
		return configdomain.Aliases{}, result.aborted(), err
	}
	selectedCommands := result.Checked()
	selectionText := DetermineAliasSelectionText(selectedCommands)
	fmt.Printf("Aliased commands: %s\n", formattedSelection(selectionText, result.aborted()))
	return DetermineAliasResult(result.CurrentSelections, allAliasableCommands, existingAliases), result.aborted(), err
}

type AliasesModel struct {
	BubbleList
	AllAliasableCommands configdomain.AliasableCommands // all Git Town commands that can be aliased
	CurrentSelections    []AliasSelection               // the status of the list entries
	OriginalAliases      configdomain.Aliases           // the Git Town aliases as they currently exist on disk
	selectedColor        termenv.Style
}

func (self AliasesModel) Checked() configdomain.AliasableCommands {
	result := configdomain.AliasableCommands{}
	for c, choice := range self.CurrentSelections {
		if choice == AliasSelectionGT {
			result = append(result, self.AllAliasableCommands[c])
		}
	}
	return result
}

func (self AliasesModel) Init() tea.Cmd {
	return nil
}

// RotateCurrentEntry switches the status of the currently selected list entry to the next status.
func (self *AliasesModel) RotateCurrentEntry() {
	var newSelection AliasSelection
	switch self.CurrentSelections[self.Cursor] {
	case AliasSelectionNone:
		commandAtCursor := self.AllAliasableCommands[self.Cursor]
		originalAlias, hasOriginalAlias := self.OriginalAliases[commandAtCursor]
		if hasOriginalAlias && originalAlias != "town "+commandAtCursor.String() {
			newSelection = AliasSelectionOther
		} else {
			newSelection = AliasSelectionGT
		}
	case AliasSelectionOther:
		newSelection = AliasSelectionGT
	case AliasSelectionGT:
		newSelection = AliasSelectionNone
	}
	self.CurrentSelections[self.Cursor] = newSelection
}

// SelectAll checks all entries in the list.
func (self *AliasesModel) SelectAll() {
	for s := range self.CurrentSelections {
		self.CurrentSelections[s] = AliasSelectionGT
	}
}

// SelectNone unchecks all entries in the list.
func (self *AliasesModel) SelectNone() {
	for s := range self.CurrentSelections {
		self.CurrentSelections[s] = AliasSelectionNone
	}
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
		self.RotateCurrentEntry()
		return self, nil
	case tea.KeyEnter:
		self.Status = dialogStatusDone
		return self, tea.Quit
	}
	switch keyMsg.String() {
	case "a":
		self.SelectAll()
	case "n":
		self.SelectNone()
	case "o":
		self.RotateCurrentEntry()
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
		s.WriteString(self.entryNumberStr(i))
		highlighted := self.Cursor == i
		selection := self.CurrentSelections[i]
		switch {
		case highlighted && selection == AliasSelectionNone:
			s.WriteString(self.Colors.selection.Styled("> [ ] " + branch))
		case highlighted && selection == AliasSelectionOther:
			s.WriteString(self.Colors.selection.Styled("> [o] " + branch))
		case highlighted && selection == AliasSelectionGT:
			s.WriteString(self.Colors.selection.Styled("> [x] " + branch))
		case !highlighted && selection == AliasSelectionNone:
			s.WriteString("  [ ] " + branch)
		case !highlighted && selection == AliasSelectionOther:
			s.WriteString("  [o] " + branch)
		case !highlighted && selection == AliasSelectionGT:
			s.WriteString(self.selectedColor.Styled("  [x] " + branch))
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

func DetermineAliasResult(selections []AliasSelection, allAliasableCommands configdomain.AliasableCommands, existingAliases configdomain.Aliases) configdomain.Aliases {
	result := configdomain.Aliases{}
	for s, selection := range selections {
		command := allAliasableCommands[s]
		switch selection {
		case AliasSelectionGT:
			result[command] = "town " + command.String()
		case AliasSelectionOther:
			result[command] = existingAliases[command]
		case AliasSelectionNone:
			// do nothing
		}
	}
	return result
}

func DetermineAliasSelectionText(selectedCommands configdomain.AliasableCommands) string {
	switch len(selectedCommands) {
	case 0:
		return "(none)"
	case len(configdomain.AllAliasableCommands()):
		return "(all)"
	default:
		return strings.Join(selectedCommands.Strings(), ", ")
	}
}

func NewAliasSelections(aliasableCommands configdomain.AliasableCommands, existingAliases configdomain.Aliases) []AliasSelection {
	result := make([]AliasSelection, len(aliasableCommands))
	for a, aliasableCommand := range aliasableCommands {
		existingAlias, exists := existingAliases[aliasableCommand]
		switch {
		case !exists:
			result[a] = AliasSelectionNone
		case existingAlias == "town "+aliasableCommand.String():
			result[a] = AliasSelectionGT
		default:
			result[a] = AliasSelectionOther
		}
	}
	return result
}

// AliasSelection encodes the status of a checkbox in the alias dialog.
type AliasSelection int

const (
	AliasSelectionNone  AliasSelection = iota // the user chose to not set this alias
	AliasSelectionGT                          // the user chose to set this alias to the corresponding Git Town command
	AliasSelectionOther                       // the user chose to keep the alias calling an external command
)
