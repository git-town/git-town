package configdomain

// AliasableCommand defines Git Town commands that can shortened via Git aliases.
type AliasableCommand string

func (self AliasableCommand) String() string { return string(self) }

type AliasableCommands []AliasableCommand

func (self AliasableCommands) Strings() []string {
	result := make([]string, len(self))
	for c, command := range self {
		result[c] = command.String()
	}
	return result
}

const (
	AliasableCommandAppend       = AliasableCommand("append")
	AliasableCommandDiffParent   = AliasableCommand("diff-parent")
	AliasableCommandHack         = AliasableCommand("hack")
	AliasableCommandKill         = AliasableCommand("kill")
	AliasableCommandPrepend      = AliasableCommand("prepend")
	AliasableCommandPropose      = AliasableCommand("propose")
	AliasableCommandRenameBranch = AliasableCommand("rename-branch")
	AliasableCommandRepo         = AliasableCommand("repo")
	AliasableCommandSetParent    = AliasableCommand("set-parent")
	AliasableCommandShip         = AliasableCommand("ship")
	AliasableCommandSync         = AliasableCommand("sync")
)

// AllAliasableCommands provides all AliasType values.
func AllAliasableCommands() AliasableCommands {
	return []AliasableCommand{
		AliasableCommandAppend,
		AliasableCommandDiffParent,
		AliasableCommandHack,
		AliasableCommandKill,
		AliasableCommandPrepend,
		AliasableCommandPropose,
		AliasableCommandRenameBranch,
		AliasableCommandRepo,
		AliasableCommandSetParent,
		AliasableCommandShip,
		AliasableCommandSync,
	}
}

func NewAliasableCommand(command string) AliasableCommand {
	for _, aliasableCommand := range AllAliasableCommands() {
		if command == aliasableCommand.String() {
			return aliasableCommand
		}
	}
	panic("unknown aliasable command: " + command)
}

func NewAliasableCommands(commands ...string) AliasableCommands {
	result := make(AliasableCommands, len(commands))
	for c, command := range commands {
		result[c] = NewAliasableCommand(command)
	}
	return result
}
