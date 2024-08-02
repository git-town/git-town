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

// provides all the keys that describe lineage entries
func (self SingleSnapshot) LineageKeys() map[LineageKey]string {
	result := map[LineageKey]string{}
	for key, value := range self {
		if lineageKey, isLineageKey := key.CheckLineage().Get(); isLineageKey {
			result[lineageKey] = value
		}
	}
	return result
}
