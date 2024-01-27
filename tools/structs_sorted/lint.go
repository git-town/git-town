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

func main() {
	folder := "."
	if len(os.Args) > 1 {
		folder = os.Args[1]
	}
	fmt.Printf("scanning folder %q\n", folder)
	issues := []string{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return err
		}
		fileSet := token.NewFileSet()
		node, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		issues = append(issues, checkFile(node, fileSet)...)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(255)
	}
	printIssues(issues)
	os.Exit(len(issues))
}

func printIssues(issues []string) {
	for _, issue := range issues {
		fmt.Println(issue)
	}
}

func checkFile(file *ast.File, fileSet *token.FileSet) []string {
	fmt.Println("checking", file.Name)
	result := []string{}
	ast.Inspect(file, func(node ast.Node) bool {
		switch typedNode := node.(type) {
		case *ast.StructType:
			result = append(result, checkStructDefinition(typedNode, fileSet)...)
		case *ast.CompositeLit:
			result = append(result, checkStructInstantiation(typedNode, fileSet)...)
		}
		return true
	})
	return result
}

func checkStructDefinition(structType *ast.StructType, fileSet *token.FileSet) []string {
	result := []string{}
	structFieldNames := structDefFieldNames(structType)
	if !sort.StringsAreSorted(structFieldNames) {
		pos := fileSet.Position(structType.Pos())
		result = append(result, fmt.Sprintf("%s:%d unsorted struct fields", pos.Filename, pos.Line))
	}
	return result
}

func checkStructInstantiation(compositeLit *ast.CompositeLit, fileSet *token.FileSet) []string {
	result := []string{}
	if _, ok := compositeLit.Type.(*ast.Ident); !ok {
		return result
	}
	structFields := structInstFieldNames(compositeLit)
	if !sort.StringsAreSorted(structFields) {
		pos := fileSet.Position(compositeLit.Pos())
		result = append(result, fmt.Sprintf("%s:%d unsorted struct fields", pos.Filename, pos.Line))
	}
	return result
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
