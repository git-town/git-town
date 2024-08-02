package configdomain

import "fmt"

// Aliases contains the Git Town releated Git aliases.
// TODO: delete this and implement this as part of the actual domain objects
type Aliases map[AliasableCommand]string

func NewAliasesFromSnapshot(snapshot SingleSnapshot) (Aliases, error) {
	result := Aliases{}
	aliasableCommands := AllAliasableCommands()
	for key, value := range snapshot.AliasKeys() {
		aliasableCommand, has := aliasableCommands.Lookup(key.String()).Get()
		if !has {
			return result, fmt.Errorf("not an aliasable command: %q", key)
		}
		result[aliasableCommand] = value
	}
	return result, nil
}
