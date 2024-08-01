package configdomain

import . "github.com/git-town/git-town/v14/src/gohacks/prelude"

type AliasableCommands []AliasableCommand

// provides the AliasableCommand matching the given Git Town command
func (self AliasableCommands) Lookup(command string) Option[AliasableCommand] {
	for _, aliasableCommand := range AllAliasableCommands() {
		if aliasableCommand.String() == command {
			return Some(aliasableCommand)
		}
	}
	return None[AliasableCommand]()
}

func (self AliasableCommands) Strings() []string {
	result := make([]string, len(self))
	for c, command := range self {
		result[c] = command.String()
	}
	return result
}

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
