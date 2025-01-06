package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// A key used for storing aliases in the Git configuration
type AliasKey Key

// tries to convert this Key into an AliasKey
func NewAliasKey(key Key) Option[AliasKey] {
	if strings.HasPrefix(key.String(), AliasKeyPrefix) {
		return Some(AliasKey(key))
	}
	return None[AliasKey]()
}

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
