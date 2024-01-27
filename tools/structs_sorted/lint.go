package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
)

func main() {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, ".", nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pack := range packs {
		for _, file := range pack.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.StructType:
					checkStructFields(node, set)
				case *ast.CompositeLit:
					if structLit, ok := node.Type.(*ast.Ident); ok {
						checkStructInstantiation(structLit, node, set)
					}
				}
				return true
			})
		}
	}
}

func checkStructFields(structType *ast.StructType, set *token.FileSet) {
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

func checkStructInstantiation(structLit *ast.Ident, compLit *ast.CompositeLit, set *token.FileSet) {
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
