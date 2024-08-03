package configdomain

// Aliases contains the Git Town releated Git aliases.
type Aliases map[AliasableCommand]string

func NewAliasesFromSnapshot(snapshot SingleSnapshot) (Aliases, error) {
	result := Aliases{}
	for key, value := range snapshot.AliasKeys() {
		result[key.AliasableCommand()] = value
	}
	return result, nil
}
