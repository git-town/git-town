package configdomain

import "github.com/git-town/git-town/v14/pkg/keys"

// SingleSnapshot contains all of the local or global Git metadata config settings.
type SingleSnapshot map[keys.Key]string

// provides all the keys that describe aliases for Git Town commands
func (self SingleSnapshot) AliasEntries() map[keys.AliasKey]string {
	result := map[keys.AliasKey]string{}
	for key, value := range self {
		if aliasKey, isAliasKey := keys.NewAliasKey(key).Get(); isAliasKey {
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
func (self SingleSnapshot) LineageEntries() map[keys.LineageKey]string {
	result := map[keys.LineageKey]string{}
	for key, value := range self {
		if lineageKey, isLineageKey := keys.NewLineageKey(key).Get(); isLineageKey {
			result[lineageKey] = value
		}
	}
	return result
}
