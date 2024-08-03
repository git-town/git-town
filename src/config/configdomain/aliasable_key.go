package configdomain

type AliasKey Key

func (self AliasKey) Key() Key {
	return Key(self)
}
