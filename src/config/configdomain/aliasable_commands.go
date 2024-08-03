package configdomain

import . "github.com/git-town/git-town/v14/src/gohacks/prelude"

type AliasableCommands []AliasableCommand

func (self AliasableCommands) CheckAliasKey(name string) Option[AliasKey] {
	for _, aliasableCommand := range self {
		keyOfCommand := aliasableCommand.Key()
		if keyOfCommand.String() == name {
			return Some(keyOfCommand)
		}
	}
	return None[AliasKey]()
}

func (self AliasableCommands) Keys() []AliasKey {
	result := make([]AliasKey, len(self))
	for a, aliasableCommand := range self {
		result[a] = aliasableCommand.Key()
	}
	return result
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
