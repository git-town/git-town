package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/invopop/jsonschema"
)

func main() {
	reflector := new(jsonschema.Reflector)
	reflector.RequiredFromJSONSchemaTags = true
	reflector.KeyNamer = camelToKebab
	schema := reflector.Reflect(&configfile.Data{}) //exhaustruct:ignore
	schema.ID = "https://www.git-town.com/git-town.toml"

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(schemaJSON))
}

// camelToKebab converts a CamelCase string to kebab-case.
// For example, "SyncStrategy" becomes "sync-strategy".
func camelToKebab(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
