package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

var (
	// file paths to ignore
	ignorePaths = []string{ //nolint:gochecknoglobals
		"vendor/",
		"tools/structs_sorted/test.go",
	}

	// struct types to ignore
	ignoreTypes = []string{ //nolint:gochecknoglobals
		"BranchSpan",
		"Change",
		"InconsistentChange",
		"Parts",
		"ProdRunner",
	}
)

func main() {
	issues := Issues{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".go") || isIgnoredPath(path) {
			return err
		}
		issues = append(issues, LintFile(path)...)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(issues) > 0 {
		fmt.Println(issues)
	}
	os.Exit(len(issues))
}

func LintFile(path string) Issues {
	result := Issues{}
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return result
	}
	ast.Inspect(file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.TypeSpec:
			result = append(result, lintStructDefinition(n, fileSet)...)
		case *ast.CompositeLit:
		}
		result = append(result, lintStructLiteral(node, fileSet)...)
		return true
	})
	return result
}

func lintStructDefinition(typeSpec *ast.TypeSpec, fileSet *token.FileSet) Issues {
	structName := typeSpec.Name.Name
	if slices.Contains(ignoreTypes, structName) {
		return Issues{}
	}
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return Issues{}
	}
	fields := structDefFieldNames(structType)
	if len(fields) == 0 {
		return Issues{}
	}
	sortedFields := make([]string, len(fields))
	copy(sortedFields, fields)
	slices.Sort(sortedFields)
	if reflect.DeepEqual(fields, sortedFields) {
		return Issues{}
	}
	return Issues{
		issue{
			expected:   sortedFields,
			position:   fileSet.Position(typeSpec.Pos()),
			structName: structName,
		},
	}
}

func lintStructLiteral(compositeLit *ast.CompositeLit, fileSet *token.FileSet) Issues {
	structType, ok := compositeLit.Type.(*ast.Ident)
	if !ok {
		return Issues{}
	}
	pos := fileSet.Position(compositeLit.Pos())
	fmt.Printf("%s:%d  %s\n", pos.Filename, pos.Line, structType.Name)
	structName := structType.Name
	if slices.Contains(ignoreTypes, structName) {
		return Issues{}
	}
	fieldNames := structInstantiationFieldNames(compositeLit)
	if len(fieldNames) == 0 {
		return Issues{}
	}
	sortedFields := make([]string, len(fieldNames))
	copy(sortedFields, fieldNames)
	slices.Sort(sortedFields)
	if reflect.DeepEqual(fieldNames, sortedFields) {
		return Issues{}
	}
	return Issues{
		issue{
			expected:   sortedFields,
			position:   fileSet.Position(compositeLit.Pos()),
			structName: structName,
		},
	}
}

func isIgnoredPath(path string) bool {
	for _, ignore := range ignorePaths {
		if strings.HasPrefix(path, ignore) {
			return true
		}
	}
	return false
}

func structDefFieldNames(structType *ast.StructType) []string {
	var result []string
	for _, field := range structType.Fields.List {
		if field.Names != nil {
			result = append(result, field.Names[0].Name)
		}
	}
	return result
}

func structInstantiationFieldNames(compositeLit *ast.CompositeLit) []string {
	result := make([]string, 0, len(compositeLit.Elts))
	for _, expr := range compositeLit.Elts {
		kvExpr, ok := expr.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kvExpr.Key.(*ast.Ident)
		if !ok {
			continue
		}
		result = append(result, ident.Name)
	}
	return result
}
