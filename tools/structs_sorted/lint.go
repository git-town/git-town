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
	ignorePaths = []string{"vendor/", "tools/structs_sorted/test.go"} //nolint:gochecknoglobals

	// struct types to ignore
	ignoreTypes = []string{"BranchSpan", "Change", "InconsistentChange", "Parts", "ProdRunner"} //nolint:gochecknoglobals
)

type issue struct {
	expected   []string
	pos        token.Position
	structName string
}

func main() {
	issues := []issue{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return err
		}
		for _, ignore := range ignorePaths {
			if strings.HasPrefix(path, ignore) {
				return nil
			}
		}
		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		issues = append(issues, checkFile(file, fileSet)...)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	printIssues(issues)
	os.Exit(len(issues))
}

func printIssues(issues []issue) {
	for _, issue := range issues {
		fmt.Printf("%s:%d:%d unsorted fields in %s. Expected order:\n\n%s\n\n", issue.pos.Filename, issue.pos.Line, issue.pos.Column, issue.structName, strings.Join(issue.expected, "\n"))
	}
}

func checkFile(file *ast.File, fileSet *token.FileSet) []issue {
	result := []issue{}
	ast.Inspect(file, func(node ast.Node) bool {
		typeSpec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}
		structName := typeSpec.Name.Name
		if slices.Contains(ignoreTypes, structName) {
			return true
		}
		fields := fieldNames(typeSpec)
		if len(fields) == 0 {
			return true
		}
		sortedFields := make([]string, len(fields))
		copy(sortedFields, fields)
		slices.Sort(sortedFields)
		if !reflect.DeepEqual(fields, sortedFields) {
			result = append(result, issue{
				pos:        fileSet.Position(node.Pos()),
				structName: structName,
				expected:   sortedFields,
			})
		}
		return true
	})
	return result
}

func fieldNames(typeSpec *ast.TypeSpec) []string {
	switch typedNode := typeSpec.Type.(type) {
	case *ast.StructType:
		return structDefFieldNames(typedNode)
	case *ast.CompositeLit:
		if _, ok := typedNode.Type.(*ast.Ident); !ok {
			return []string{}
		}
		return structInstFieldNames(typedNode)
	}
	return []string{}
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

func structInstFieldNames(compositeLit *ast.CompositeLit) []string {
	var result []string
	for _, expr := range compositeLit.Elts {
		if kvExpr, ok := expr.(*ast.KeyValueExpr); ok {
			if ident, ok := kvExpr.Key.(*ast.Ident); ok {
				result = append(result, ident.Name)
			}
		}
	}
	return result
}
