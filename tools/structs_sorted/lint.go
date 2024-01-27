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
				if structType, ok := n.(*ast.StructType); ok {
					checkStructFields(structType, set)
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
		fmt.Printf("Struct fields are not in alphabetical order in file %s at line %d\n",
			set.Position(structType.Pos()).Filename,
			set.Position(structType.Pos()).Line)
	}
}
