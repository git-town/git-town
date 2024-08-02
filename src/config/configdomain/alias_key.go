package configdomain

// a Key that contains an alias for a Git Town command
type AliasKey Key

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
