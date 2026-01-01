package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/invopop/jsonschema"
)

func main() {
	reflector := new(jsonschema.Reflector)
	reflector.RequiredFromJSONSchemaTags = true
	schema := reflector.Reflect(&configfile.Data{}) //exhaustruct:ignore
	schema.ID = "https://www.git-town.com/git-town.toml"

	// Post-process the schema to use TOML tag names instead of Go field names
	convertToTomlNames(schema, reflect.TypeOf((*configfile.Data)(nil)).Elem())

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(schemaJSON))
}

// convertToTomlNames recursively walks through the schema and renames properties
// based on the TOML struct tags
func convertToTomlNames(schema *jsonschema.Schema, t reflect.Type) {
	if schema.Definitions != nil {
		// Collect and sort definition names for deterministic iteration
		defNames := slices.Collect(maps.Keys(schema.Definitions))
		sort.Strings(defNames)

		for _, defName := range defNames {
			defSchema := schema.Definitions[defName]
			// Find the corresponding type
			var defType reflect.Type
			switch defName {
			case "Branches":
				defType = reflect.TypeOf((*configfile.Branches)(nil)).Elem()
			case "Create":
				defType = reflect.TypeOf((*configfile.Create)(nil)).Elem()
			case "Hosting":
				defType = reflect.TypeOf((*configfile.Hosting)(nil)).Elem()
			case "Propose":
				defType = reflect.TypeOf((*configfile.Propose)(nil)).Elem()
			case "Ship":
				defType = reflect.TypeOf((*configfile.Ship)(nil)).Elem()
			case "Sync":
				defType = reflect.TypeOf((*configfile.Sync)(nil)).Elem()
			case "SyncStrategy":
				defType = reflect.TypeOf((*configfile.SyncStrategy)(nil)).Elem()
			case "Data":
				defType = reflect.TypeOf((*configfile.Data)(nil)).Elem()
			}
			if defType != nil {
				renameProperties(defSchema, defType)
			}
		}
	}

	// Also rename properties in the root schema
	renameProperties(schema, t)
}

// renameProperties renames the properties in a schema based on TOML struct tags
func renameProperties(schema *jsonschema.Schema, t reflect.Type) {
	if schema.Properties == nil {
		return
	}

	newProperties := make(map[string]*jsonschema.Schema)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tomlTag := field.Tag.Get("toml")
		if tomlTag == "" {
			continue
		}

		// Extract the field name from the tag (before any comma-separated options)
		parts := strings.Split(tomlTag, ",")
		tomlName := parts[0]

		// Find the property with the Go field name and rename it
		if prop, exists := schema.Properties.Get(field.Name); exists {
			newProperties[tomlName] = prop
		}
	}

	// Replace the properties with the renamed ones
	// Sort keys for deterministic iteration
	keys := make([]string, 0, len(newProperties))
	for k := range newProperties { // okay to iterate the map in random order
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := newProperties[k]
		schema.Properties.Set(k, v)
		// Remove the old property if it had a different name
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tomlTag := field.Tag.Get("toml")
			if tomlTag != "" {
				parts := strings.Split(tomlTag, ",")
				tomlName := parts[0]
				if tomlName == k && field.Name != k {
					schema.Properties.Delete(field.Name)
				}
			}
		}
	}
}
