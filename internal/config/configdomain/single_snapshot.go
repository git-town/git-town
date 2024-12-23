package configdomain

// SingleSnapshot contains all of the local or global Git metadata config settings.
type SingleSnapshot map[Key]string

// provides all the keys that describe aliases for Git Town commands
func (self SingleSnapshot) AliasEntries() map[AliasKey]string {
	result := map[AliasKey]string{}
	for key, value := range self {
		if aliasKey, isAliasKey := NewAliasKey(key).Get(); isAliasKey {
			result[aliasKey] = value
		}
	}
	return result
}

func (self SingleSnapshot) Aliases() Aliases {
	aliasEntries := self.AliasEntries()
	result := make(Aliases, len(aliasEntries))
	for key, value := range aliasEntries {
		result[key.AliasableCommand()] = value
	}
	return result
}

// provides all the keys that describe lineage entries
func (self SingleSnapshot) LineageEntries() map[LineageKey]string {
	result := map[LineageKey]string{}
	for key, value := range self {
		if lineageKey, isLineageKey := ParseLineageKey(key).Get(); isLineageKey {
			result[lineageKey] = value
		}
	}
	return result
}

// provides all the keys that describe branch type overrides
func (self SingleSnapshot) BranchTypeOverrideEntries() map[BranchTypeOverrideKey]string {
	result := map[BranchTypeOverrideKey]string{}
	for key, value := range self {
		if branchTypeKey, isBranchTypeKey := ParseBranchTypeOverrideKey(key).Get(); isBranchTypeKey {
			result[branchTypeKey] = value
		}
	}
	return result
}
