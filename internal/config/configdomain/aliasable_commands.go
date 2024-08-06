package configdomain

import . "github.com/git-town/git-town/v14/pkg/prelude"

type AliasableCommands []AliasableCommand

// provides the AliasKey matching the given key name
func (self AliasableCommands) LookupKey(name string) Option[AliasKey] {
	for _, aliasableCommand := range self {
		keyOfCommand := aliasableCommand.Key()
		if keyOfCommand.String() == name {
			return Some(keyOfCommand)
		}
	}
	return None[AliasKey]()
}

func (self AliasableCommands) Strings() []string {
	result := make([]string, len(self))
	for c, command := range self {
		result[c] = command.String()
	}
	return result
}
