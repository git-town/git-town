package configdomain

import "fmt"

// Aliases contains the Git Town releated Git aliases.
type Aliases map[AliasableCommand]string

func NewAliasesFromSnapshot(snapshot SingleSnapshot) (Aliases, error) {
	result := Aliases{}
	for key, value := range snapshot.AliasKeys() {
		aliasableCommand, has := AliasableCommandForKey(key).Get()
		if !has {
			return result, fmt.Errorf("not an aliasable command: %q", key)
		}
		result[aliasableCommand] = value
	}
	return result, nil
}
