package configdomain

// A key used for storing aliases in the Git configuration
type AliasKey Key

func (self AliasKey) Key() Key {
	return Key(self)
}

func (self AliasKey) String() string {
	return string(self)
}
