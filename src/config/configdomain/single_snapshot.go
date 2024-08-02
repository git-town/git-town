package configdomain

// SingleSnapshot contains all of the local or global Git metadata config settings.
type SingleSnapshot map[Key]string

// provides all the keys that describe aliases for Git Town commands
func (self SingleSnapshot) AliasKeys() map[Key]string {
	result := map[Key]string{}
	for key, value := range self {
		if key.IsAliasKey() {
			result[key] = value
		}
	}
	return result
}

// provides all the keys that describe lineage entries
func (self SingleSnapshot) LineageKeys() map[Key]string {
	result := map[Key]string{}
	for key, value := range self {
		if key.IsLineage() {
			result[key] = value
		}
	}
	return result
}
