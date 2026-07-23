package configdomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/mapstools"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// SingleSnapshot contains all of the local or global Git metadata config settings.
type SingleSnapshot map[Key]string

// AliasEntries provides all the keys that describe aliases for Git Town commands
func (self SingleSnapshot) AliasEntries() map[AliasKey]string {
	result := map[AliasKey]string{}
	for key, value := range self { // okay to iterate the map in random order here because we assign to a new map
		if aliasKey, isAliasKey := NewAliasKey(key).Get(); isAliasKey {
			result[aliasKey] = value
		}
	}
	return result
}

func (self SingleSnapshot) Aliases() Aliases {
	aliasEntries := self.AliasEntries()
	result := make(Aliases, len(aliasEntries))
	for key, value := range mapstools.SortedKeyValues(aliasEntries) {
		result[key.AliasableCommand()] = value
	}
	return result
}

// GetOpt provides the snapshot content for the given key as an Option.
func (self SingleSnapshot) GetOpt(key Key) Option[string] {
	if value, has := self[key]; has {
		return Some(value)
	}
	return None[string]()
}

// LineageEntries provides all the keys that describe lineage entries
func (self SingleSnapshot) LineageEntries() map[LineageKey]string {
	result := map[LineageKey]string{}
	for key, value := range self { // okay to iterate the map in random order here because we assign to a new map
		if lineageKey, isLineageKey := ParseLineageKey(key).Get(); isLineageKey {
			result[lineageKey] = value
		}
	}
	return result
}
