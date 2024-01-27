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

func sortStructFields(structType *ast.StructType) {
	sort.Slice(structType.Fields.List, func(a, b int) bool {
		return structType.Fields.List[a].Names[0].Name < structType.Fields.List[b].Names[0].Name
	})
}

func sortCompositeLitFields(x *ast.CompositeLit) {
	sort.Slice(x.Elts, func(i, j int) bool {
		return x.Elts[i].(*ast.KeyValueExpr).Key.(*ast.Ident).Name < x.Elts[j].(*ast.KeyValueExpr).Key.(*ast.Ident).Name
	})
}
