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
	filepath.Walk(".", func(filePath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(filePath) != ".go" {
			return err
		}
		fileSet := token.NewFileSet()
		astFile, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		ast.Inspect(astFile, func(node ast.Node) bool {
			switch nodeType := node.(type) {
			case *ast.StructType:
				sortStructFields(nodeType)
			case *ast.CompositeLit:
				sortCompositeLitFields(nodeType)
			}
			return true
		})
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		return format.Node(file, fileSet, astFile)
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
