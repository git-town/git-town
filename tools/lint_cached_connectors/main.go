package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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
	// Define connector pairs to check
	connectorPairs := []ConnectorPair{
		{
			Package:      "bitbucketcloud",
			UncachedFile: "internal/forge/bitbucketcloud/api_connector.go",
			UncachedType: "APIConnector",
			CachedFile:   "internal/forge/bitbucketcloud/cached_api_connector.go",
			CachedType:   "CachedAPIConnector",
		},
		{
			Package:      "bitbucketdatacenter",
			UncachedFile: "internal/forge/bitbucketdatacenter/api_connector.go",
			UncachedType: "APIConnector",
			CachedFile:   "internal/forge/bitbucketdatacenter/cached_api_connector.go",
			CachedType:   "CachedAPIConnector",
		},
		{
			Package:      "forgejo",
			UncachedFile: "internal/forge/forgejo/api_connector.go",
			UncachedType: "APIConnector",
			CachedFile:   "internal/forge/forgejo/cached_api_connector.go",
			CachedType:   "CachedAPIConnector",
		},
		{
			Package:      "gitea",
			UncachedFile: "internal/forge/gitea/api_connector.go",
			UncachedType: "APIConnector",
			CachedFile:   "internal/forge/gitea/cached_api_connector.go",
			CachedType:   "CachedAPIConnector",
		},
		{
			Package:      "github",
			UncachedFile: "internal/forge/github/api_connector.go",
			UncachedType: "APIConnector",
			CachedFile:   "internal/forge/github/cached_api_connector.go",
			CachedType:   "CachedAPIConnector",
		},
		{
			Package:      "gitlab",
			UncachedFile: "internal/forge/gitlab/api_connector.go",
			UncachedType: "APIConnector",
			CachedFile:   "internal/forge/gitlab/cached_api_connector.go",
			CachedType:   "CachedAPIConnector",
		},
		{
			Package:      "gh",
			UncachedFile: "internal/forge/gh/connector.go",
			UncachedType: "Connector",
			CachedFile:   "internal/forge/gh/cached_connector.go",
			CachedType:   "CachedConnector",
		},
		{
			Package:      "glab",
			UncachedFile: "internal/forge/glab/connector.go",
			UncachedType: "Connector",
			CachedFile:   "internal/forge/glab/cached_connector.go",
			CachedType:   "CachedConnector",
		},
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
