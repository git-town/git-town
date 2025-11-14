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

	"github.com/git-town/git-town/v22/pkg/set"
)

// ConnectorPair represents a pair of cached and uncached connectors
type ConnectorPair struct {
	CachedFile   string // file path of the cached connector
	CachedType   string // type name of the cached connector
	Package      string // package name of the connector pair
	UncachedFile string // file path of the uncached connector
	UncachedType string // type name of the uncached connector
}

// InterfaceImplementation represents a type implementing an interface
type InterfaceImplementation struct {
	InterfaceName string         // name of the interface that is implemented
	Position      token.Position // file, line, and column where the interface is implemented
	TypeName      string         // type of the connector that implements the interface
}

func main() {
	if err := check(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func check() error {
	connectorPairs, err := connectorPairs("internal/forge")
	if err != nil {
		return fmt.Errorf("Error discovering connector pairs: %v\n", err)
	}
	if len(connectorPairs) == 0 {
		return fmt.Errorf("No connector pairs found")
	}
	for _, pair := range connectorPairs {
		uncachedInterfaces, err := implementedInterfaces(pair.UncachedFile, pair.UncachedType)
		if err != nil {
			return fmt.Errorf("Error parsing %s: %v", pair.UncachedFile, err)
		}
		cachedInterfaces, err := implementedInterfaces(pair.CachedFile, pair.CachedType)
		if err != nil {
			return fmt.Errorf("Error parsing %s: %v", pair.CachedFile, err)
		}
		for _, missing := range missingInterfaces(uncachedInterfaces, cachedInterfaces) {
			return fmt.Errorf(
				"%s:1  %s needs to implement %s\nbecause %s implements it in %s:%d",
				pair.CachedFile,
				pair.CachedType,
				missing.InterfaceName,
				pair.UncachedType,
				missing.Position.Filename,
				missing.Position.Line,
			)
		}
	}
	return nil
}

// connectorPairs discovers all cached/uncached connector pairs in the given directory
func connectorPairs(dir string) ([]ConnectorPair, error) {
	var result []ConnectorPair
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkgName := entry.Name()
		pkgPath := filepath.Join(dir, pkgName)
		globPattern := filepath.Join(pkgPath, "cached_*.go")
		cachedFiles, err := filepath.Glob(globPattern)
		if err != nil {
			return nil, fmt.Errorf("globbing %s in %s: %w", globPattern, pkgPath, err)
		}
		for _, cachedFile := range cachedFiles {
			uncachedFile := uncachedFilePath(cachedFile, pkgPath)
			if _, err := os.Stat(uncachedFile); os.IsNotExist(err) {
				return nil, fmt.Errorf("uncached file %s does not exist", uncachedFile)
			}
			cachedType, err := primaryTypeName(cachedFile)
			if err != nil {
				return nil, fmt.Errorf("extracting type from %s: %w", cachedFile, err)
			}
			uncachedType, err := primaryTypeName(uncachedFile)
			if err != nil {
				return nil, fmt.Errorf("extracting type from %s: %w", uncachedFile, err)
			}
			result = append(result, ConnectorPair{
				CachedFile:   cachedFile,
				CachedType:   cachedType,
				Package:      pkgName,
				UncachedFile: uncachedFile,
				UncachedType: uncachedType,
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Package < result[j].Package
	})
	return result, nil
}

// uncachedFilePath extracts the uncached file path from a cached file path.
// "cached_api_connector.go" -> "api_connector.go"
// "cached_connector.go" -> "connector.go"
func uncachedFilePath(cachedFile, pkgPath string) string {
	baseName := filepath.Base(cachedFile)
	uncachedName := strings.TrimPrefix(baseName, "cached_")
	return filepath.Join(pkgPath, uncachedName)
}

// primaryTypeName parses a Go file and extracts the primary struct type name.
// It looks for the first struct type declaration in the file.
func primaryTypeName(filePath string) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filePath, nil, 0)
	if err != nil {
		return "", err
	}
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
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				return typeSpec.Name.Name, nil
			}
		}
	}
	return "", fmt.Errorf("no struct type found in %s", filePath)
}

// implementedInterfaces parses a Go file and extracts all interfaces
// that the given type name implements
func implementedInterfaces(filePath, typeName string) ([]InterfaceImplementation, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var implementations []InterfaceImplementation

	// First, find all variables of our target type
	// Pattern: var myConnector MyConnectorType
	typeVars := set.New[string]()
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
			if getTypeName(valueSpec.Type) == typeName {
				for _, name := range valueSpec.Names {
					typeVars.Add(name.Name)
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
				valueName := valueVarName(valueSpec.Values[0])
				if typeVars.Contains(valueName) {
					position := fileSet.Position(valueSpec.Pos())
					implementations = append(implementations, InterfaceImplementation{
						Position:      position,
						InterfaceName: interfaceName,
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

// valueVarName extracts the variable name from an expression, handling both
// direct references (varName) and pointer references (&varName)
func valueVarName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.UnaryExpr:
		// Handle &varName
		if t.Op == token.AND {
			return valueVarName(t.X)
		}
		return ""
	default:
		return ""
	}
}

// missingInterfaces returns interfaces that are in 'expected' but not in 'actual'
func missingInterfaces(expected, actual []InterfaceImplementation) []InterfaceImplementation {
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
	sort.Slice(missing, func(i, j int) bool {
		return strings.Compare(missing[i].InterfaceName, missing[j].InterfaceName) < 0
	})
	return missing
}
