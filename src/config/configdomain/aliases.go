package configdomain

// Aliases contains the Git Town releated Git aliases.
// TODO: delete this and implement this as part of the actual domain objects
type Aliases map[AliasableCommand]string

func NewAliasesFromSnapshot(snapshot SingleSnapshot) (Aliases, error) {
	result := Aliases{}
	aliasableCommands := AllAliasableCommands()
	for key, value := range snapshot.AliasKeys() {
		result[key.AliasableCommand()] = value
	}
	return result, nil
}
