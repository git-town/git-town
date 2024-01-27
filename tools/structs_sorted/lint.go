package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	filepath.Walk("src", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.StructType:
				sortStructFields(x)
			case *ast.CompositeLit:
				sortCompositeLitFields(x)
			}
			return true
		})

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := format.Node(file, fset, node); err != nil {
			return err
		}

		return nil
	})
}

func sortStructFields(x *ast.StructType) {
	sort.Slice(x.Fields.List, func(i, j int) bool {
		return x.Fields.List[i].Names[0].Name < x.Fields.List[j].Names[0].Name
	})
}

func sortCompositeLitFields(x *ast.CompositeLit) {
	sort.Slice(x.Elts, func(i, j int) bool {
		return x.Elts[i].(*ast.KeyValueExpr).Key.(*ast.Ident).Name < x.Elts[j].(*ast.KeyValueExpr).Key.(*ast.Ident).Name
	})
}
