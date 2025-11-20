package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/invopop/jsonschema"
)

func main() {
	reflector := new(jsonschema.Reflector)
	schema := reflector.Reflect(&configfile.Data{})
	schema.ID = "https://www.git-town.com/git-town.toml"
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(schemaJSON))
}
