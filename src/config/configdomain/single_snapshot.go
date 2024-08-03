package configdomain

// SingleSnapshot contains all of the local or global Git metadata config settings.
type SingleSnapshot map[Key]string

func (self SingleSnapshot) Aliases() Aliases {
	aliasEntries := self.AliasEntries()
	result := make(Aliases, len(aliasEntries))
	for key, value := range aliasEntries {
		result[key.AliasableCommand()] = value
	}
	return result
}

// provides all the keys that describe aliases for Git Town commands
func (self SingleSnapshot) AliasEntries() map[AliasKey]string {
	result := map[AliasKey]string{}
	for key, value := range self {
		if aliasKey, isAliasKey := key.ToAliasKey().Get(); isAliasKey {
			result[aliasKey] = value
		}
	}
	return result
}

// provides all the keys that describe lineage entries
func (self SingleSnapshot) LineageEntries() map[LineageKey]string {
	result := map[LineageKey]string{}
	for key, value := range self {
		if lineageKey, isLineageKey := key.CheckLineage().Get(); isLineageKey {
			result[lineageKey] = value
		}
	}
	return result
}
