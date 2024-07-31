package configdomain

import "fmt"

// Aliases contains the Git Town releated Git aliases.
type Aliases map[AliasableCommand]string

func NewAliasesFromSnapshot(snapshot SingleSnapshot) (Aliases, error) {
	result := Aliases{}
	for key, value := range snapshot {
		if key.IsAliasKey() {
			aliasableCommand, has := AliasableCommandForKey(key).Get()
			if !has {
				return result, fmt.Errorf("not an aliasable command: " + key.String())
			}
			result[aliasableCommand] = value
		}
	}
	return result, nil
}
