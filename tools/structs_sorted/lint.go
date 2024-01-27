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

func checkStructDefinition(structType *ast.StructType, set *token.FileSet) {
	var fieldNames []string
	for _, field := range structType.Fields.List {
		if field.Names != nil {
			fieldNames = append(fieldNames, field.Names[0].Name)
		}
	}

	if !sort.StringsAreSorted(fieldNames) {
		fmt.Printf("%s:%d unsorted struct fields\n",
			set.Position(structType.Pos()).Filename,
			set.Position(structType.Pos()).Line)
	}
}

func checkStructInstantiation(compLit *ast.CompositeLit, set *token.FileSet) {
	_, ok := compLit.Type.(*ast.Ident)
	if !ok {
		return
	}
	var fieldNames []string
	for _, expr := range compLit.Elts {
		if kvExpr, ok := expr.(*ast.KeyValueExpr); ok {
			if ident, ok := kvExpr.Key.(*ast.Ident); ok {
				fieldNames = append(fieldNames, ident.Name)
			}
		}
	}

	if !sort.StringsAreSorted(fieldNames) {
		fmt.Printf("%s:%d unsorted struct fields\n",
			set.Position(compLit.Pos()).Filename,
			set.Position(compLit.Pos()).Line,
		)
	}
}
