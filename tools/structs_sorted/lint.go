package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
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

		sort.Slice(structType.Fields.List, func(a, b int) bool {
			if structType.Fields.List[a].Names == nil || structType.Fields.List[b].Names == nil {
				return false
			}
			return structType.Fields.List[a].Names[0].Name < structType.Fields.List[b].Names[0].Name
		})

		return true
	})

	file, err := os.Create("test_sorted.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	err = format.Node(file, fset, node)
	if err != nil {
		fmt.Println(err)
		return
	}
}
