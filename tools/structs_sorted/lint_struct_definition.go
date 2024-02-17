package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"slices"
)

func lintStructDefinition(typeSpec *ast.TypeSpec, fileSet *token.FileSet, issues *Issues) {
	structName := typeSpec.Name.Name
	if slices.Contains(ignoreTypes, structName) {
		return
	}
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return
	}
	fields := structDefFieldNames(structType)
	if len(fields) == 0 {
		return
	}
	sortedFields := make([]string, len(fields))
	copy(sortedFields, fields)
	slices.Sort(sortedFields)
	if reflect.DeepEqual(fields, sortedFields) {
		return
	}
	*issues = append(*issues, issue{
		expected: sortedFields,
		position: fileSet.Position(typeSpec.Pos()),
	})
}

func structDefFieldNames(structType *ast.StructType) []string {
	var result []string
	for _, field := range structType.Fields.List {
		if field.Names != nil {
			result = append(result, field.Names[0].Name)
		}
	}
	return result
}
