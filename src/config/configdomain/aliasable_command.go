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
	AliasableCommandCompress     = AliasableCommand("compress")
	AliasableCommandContribute   = AliasableCommand("contribute")
	AliasableCommandDiffParent   = AliasableCommand("diff-parent")
	AliasableCommandHack         = AliasableCommand("hack")
	AliasableCommandKill         = AliasableCommand("kill")
	AliasableCommandObserve      = AliasableCommand("observe")
	AliasableCommandPark         = AliasableCommand("park")
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
		AliasableCommandCompress,
		AliasableCommandContribute,
		AliasableCommandDiffParent,
		AliasableCommandHack,
		AliasableCommandKill,
		AliasableCommandObserve,
		AliasableCommandPark,
		AliasableCommandPrepend,
		AliasableCommandPropose,
		AliasableCommandRenameBranch,
		AliasableCommandRepo,
		AliasableCommandSetParent,
		AliasableCommandShip,
		AliasableCommandSync,
	}
}
