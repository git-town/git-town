package configdomain

import "strings"

// A key used for storing aliases in the Git configuration
type AliasKey Key

// provides the AliasableCommand matching this AliasKey
func (self AliasKey) AliasableCommand() AliasableCommand {
	commandName := strings.TrimPrefix(self.String(), AliasKeyPrefix)
	return AliasableCommand(commandName)
}

// provides the generic Key that this AliasKey represents
func (self AliasKey) Key() Key {
	return Key(self)
}

func (self AliasKey) String() string {
	return string(self)
}

const AliasKeyPrefix = "alias."
