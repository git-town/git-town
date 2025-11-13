package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ConnectorPair represents a pair of cached and uncached connectors
type ConnectorPair struct {
	CachedFile   string
	CachedType   string
	Package      string
	UncachedFile string
	UncachedType string
}

// InterfaceImplementation represents a type implementing an interface
type InterfaceImplementation struct {
	FilePath      string
	InterfaceName string
	LineNumber    int
	TypeName      string
}

func main() {
	// Discover connector pairs dynamically
	connectorPairs, err := discoverConnectorPairs("internal/forge")
	if err != nil {
		fmt.Printf("Error discovering connector pairs: %v\n", err)
		os.Exit(1)
	}

	if len(connectorPairs) == 0 {
		fmt.Println("No connector pairs found")
		os.Exit(1)
	}

	var allErrors []string

	for _, pair := range connectorPairs {
		uncachedInterfaces, err := extractInterfaceImplementations(pair.UncachedFile, pair.UncachedType)
		if err != nil {
			allErrors = append(allErrors, fmt.Sprintf("Error parsing %s: %v", pair.UncachedFile, err))
			continue
		}

		cachedInterfaces, err := extractInterfaceImplementations(pair.CachedFile, pair.CachedType)
		if err != nil {
			allErrors = append(allErrors, fmt.Sprintf("Error parsing %s: %v", pair.CachedFile, err))
			continue
		}

		// Check if cached connector implements all interfaces that uncached does
		missing := findMissingInterfaces(uncachedInterfaces, cachedInterfaces)
		for _, iface := range missing {
			allErrors = append(allErrors, fmt.Sprintf(
				"%s:%d: %s does not implement interface %s (implemented by %s in %s:%d)",
				pair.CachedFile,
				1, // Could be improved to show actual line
				pair.CachedType,
				iface.InterfaceName,
				pair.UncachedType,
				pair.UncachedFile,
				iface.LineNumber,
			))
		}
	}

	if len(allErrors) > 0 {
		for _, err := range allErrors {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

// discoverConnectorPairs scans the forge directory and dynamically discovers
// all cached/uncached connector pairs
func discoverConnectorPairs(forgeDir string) ([]ConnectorPair, error) {
	var pairs []ConnectorPair

	// Walk through all subdirectories in the forge directory
	entries, err := os.ReadDir(forgeDir)
	if err != nil {
		return nil, fmt.Errorf("reading forge directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pkgName := entry.Name()
		pkgPath := filepath.Join(forgeDir, pkgName)

		// Look for cached_*.go files
		cachedFiles, err := filepath.Glob(filepath.Join(pkgPath, "cached_*.go"))
		if err != nil {
			return nil, fmt.Errorf("globbing cached files in %s: %w", pkgPath, err)
		}

		for _, cachedFile := range cachedFiles {
			// Extract base name from cached file
			// e.g., "cached_api_connector.go" -> "api_connector.go"
			//       "cached_connector.go" -> "connector.go"
			baseName := filepath.Base(cachedFile)
			uncachedName := strings.TrimPrefix(baseName, "cached_")
			uncachedFile := filepath.Join(pkgPath, uncachedName)

			// Check if the uncached file exists
			if _, err := os.Stat(uncachedFile); os.IsNotExist(err) {
				continue
			}

			// Extract type names from the files
			cachedType, err := extractPrimaryTypeName(cachedFile)
			if err != nil {
				return nil, fmt.Errorf("extracting cached type from %s: %w", cachedFile, err)
			}

			uncachedType, err := extractPrimaryTypeName(uncachedFile)
			if err != nil {
				return nil, fmt.Errorf("extracting uncached type from %s: %w", uncachedFile, err)
			}

			if cachedType != "" && uncachedType != "" {
				pairs = append(pairs, ConnectorPair{
					CachedFile:   cachedFile,
					CachedType:   cachedType,
					Package:      pkgName,
					UncachedFile: uncachedFile,
					UncachedType: uncachedType,
				})
			}
		}
	}

	// Sort pairs by package name for consistent output
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Package < pairs[j].Package
	})

	return pairs, nil
}

// extractPrimaryTypeName parses a Go file and extracts the primary struct type name.
// It looks for the first struct type declaration in the file.
func extractPrimaryTypeName(filePath string) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filePath, nil, 0)
	if err != nil {
		return "", err
	}

	// Look for struct type declarations
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// Check if it's a struct type
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				return typeSpec.Name.Name, nil
			}
		}
	}

	return "", fmt.Errorf("no struct type found in %s", filePath)
}

// extractInterfaceImplementations parses a Go file and extracts all interface implementations
// for the given type name
func extractInterfaceImplementations(filePath, typeName string) ([]InterfaceImplementation, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var implementations []InterfaceImplementation

	// First, find all variables of our target type
	// Pattern: var myConnector MyConnectorType
	typeVars := make(map[string]bool) // variable name -> true
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || valueSpec.Type == nil {
				continue
			}

			// Check if the type matches our target type
			if getTypeName(valueSpec.Type) == typeName {
				for _, name := range valueSpec.Names {
					typeVars[name.Name] = true
				}
			}
		}
	}

	// Now find interface type checks
	// Pattern: var _ InterfaceName = myConnectorVar or &myConnectorVar
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// Check if this is a type assertion (var _ Interface = variable)
			if len(valueSpec.Names) == 1 && valueSpec.Names[0].Name == "_" && valueSpec.Type != nil && len(valueSpec.Values) > 0 {
				// Get the interface name
				interfaceName := getTypeName(valueSpec.Type)

				// Check if the value is a variable of our target type
				// Handle both: varName and &varName
				valueName := getValueVarName(valueSpec.Values[0])
				if typeVars[valueName] {
					position := fileSet.Position(valueSpec.Pos())
					implementations = append(implementations, InterfaceImplementation{
						FilePath:      filePath,
						InterfaceName: interfaceName,
						LineNumber:    position.Line,
						TypeName:      typeName,
					})
				}
			}
		}
	}

	return implementations, nil
}

// getTypeName extracts the type name from an AST expression
func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		// For qualified names like forgedomain.Connector
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
		return t.Sel.Name
	case *ast.CompositeLit:
		return getTypeName(t.Type)
	default:
		return ""
	}
}

// getValueVarName extracts the variable name from an expression, handling both
// direct references (varName) and pointer references (&varName)
func getValueVarName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.UnaryExpr:
		// Handle &varName
		if t.Op == token.AND {
			return getValueVarName(t.X)
		}
		return ""
	default:
		return ""
	}
}

// findMissingInterfaces returns interfaces that are in 'expected' but not in 'actual'
func findMissingInterfaces(expected, actual []InterfaceImplementation) []InterfaceImplementation {
	actualSet := make(map[string]bool)
	for _, impl := range actual {
		actualSet[impl.InterfaceName] = true
	}

	var missing []InterfaceImplementation
	for _, impl := range expected {
		if !actualSet[impl.InterfaceName] {
			missing = append(missing, impl)
		}
	}

	// Sort for consistent output
	sort.Slice(missing, func(i, j int) bool {
		return strings.Compare(missing[i].InterfaceName, missing[j].InterfaceName) < 0
	})

	return missing
}
