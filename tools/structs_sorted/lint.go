package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Type == nil {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok || structType.Fields == nil {
			return true
		}

		fields := make([]string, len(structType.Fields.List))
		for f, field := range structType.Fields.List {
			if field.Names != nil {
				fields[f] = field.Names[0].Name
			}
		}

		if !sort.StringsAreSorted(fields) {
			fmt.Printf("Struct %s has unsorted fields\n", typeSpec.Name.Name)
		}

		return true
	})
}
