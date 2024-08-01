package configdomain

// SingleSnapshot contains all of the local or global Git metadata config settings.
type SingleSnapshot map[Key]string

func (self SingleSnapshot) AliasKeys() map[AliasKey]string {
	result := map[AliasKey]string{}
	for key, value := range self {
		if key.IsAliasKey() {
			result[AliasKey(key)] = value
		}
	}
	return result
}

func (self SingleSnapshot) LineageKeys() map[LineageKey]string {
	result := map[LineageKey]string{}
	for key, value := range self {
		if key.IsLineage() {
			result[LineageKey(key)] = value
		}
	}
	return result
}
