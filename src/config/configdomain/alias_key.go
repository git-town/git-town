package configdomain

import . "github.com/git-town/git-town/v14/src/gohacks/prelude"

// a Key that contains an alias for a Git Town command
type AliasKey Key

func (self AliasKey) AliasableCommand() Option[AliasableCommand] {
	selfKey := self.Key()
	for _, aliasableCommand := range AllAliasableCommands() {
		if KeyForAliasableCommand(aliasableCommand) == selfKey {
			return Some(aliasableCommand)
		}
	}
	return None[AliasableCommand]()
}

func (self AliasKey) Key() Key {
	return Key(self)
}

func (self AliasKey) String() string {
	return string(self)
}

func NewAliasKey(command string) (AliasKey, error) {
	for _, aliasableCommand := range AllAliasableCommands() {
		aliasableCommand
	}
}
