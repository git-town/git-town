package configdomain

type AliasableCommands []AliasableCommand

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
