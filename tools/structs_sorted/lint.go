package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"sort"
)

func main() {
	fileSet := token.NewFileSet()
	packs, err := parser.ParseDir(fileSet, ".", nil, 0)
	if err != nil {
		log.Fatalln(err)
	}
	for _, pack := range packs {
		for _, file := range pack.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				switch typedNode := node.(type) {
				case *ast.StructType:
					checkStructDefinition(typedNode, fileSet)
				case *ast.CompositeLit:
					checkStructInstantiation(typedNode, fileSet)
				}
				return true
			})
		}
	}
}

func checkStructDefinition(structType *ast.StructType, fileSet *token.FileSet) {
	var fieldNames []string
	for _, field := range structType.Fields.List {
		if field.Names != nil {
			fieldNames = append(fieldNames, field.Names[0].Name)
		}
	}
	if !sort.StringsAreSorted(fieldNames) {
		pos := fileSet.Position(structType.Pos())
		fmt.Printf("%s:%d unsorted struct fields\n", pos.Filename, pos.Line)
	}
}

func checkStructInstantiation(compositeLit *ast.CompositeLit, fileSet *token.FileSet) {
	if _, ok := compositeLit.Type.(*ast.Ident); !ok {
		return
	}
	structFields := structInstantiationFields(compositeLit)
	if !sort.StringsAreSorted(structFields) {
		pos := fileSet.Position(compositeLit.Pos())
		fmt.Printf("%s:%d unsorted struct fields\n", pos.Filename, pos.Line)
	}
}

func structInstantiationFields(compositeLit *ast.CompositeLit) []string {
	var fieldNames []string
	for _, expr := range compositeLit.Elts {
		if kvExpr, ok := expr.(*ast.KeyValueExpr); ok {
			if ident, ok := kvExpr.Key.(*ast.Ident); ok {
				fieldNames = append(fieldNames, ident.Name)
			}
		}
	}
	return fieldNames
}
