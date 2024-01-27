package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
)

func processFile(path string, fileSet *token.FileSet) error {
	astFile, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.Inspect(astFile, func(astNode ast.Node) bool {
		typeSpec, ok := astNode.(*ast.TypeSpec)
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
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return format.Node(file, fileSet, astFile)
}

func main() {
	fileSet := token.NewFileSet()
	err := filepath.Walk("src", func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil || fileInfo.IsDir() || filepath.Ext(path) != ".go" {
			return err
		}
		return processFile(path, fileSet)
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
