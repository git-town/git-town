package configdomain

// SingleSnapshot contains all of the local or global Git metadata config settings.
type SingleSnapshot map[Key]string

func (self SingleSnapshot) AliasKeys() map[Key]string {
	result := map[Key]string{}
	for key, value := range self {
		if key.IsAliasKey() {
			result[key] = value
		}
	}
	return result
}

func (self SingleSnapshot) LineageKeys() map[Key]string {
	result := map[Key]string{}
	for key, value := range self {
		if key.IsLineage() {
			result[key] = value
		}
	}
	return result
}
