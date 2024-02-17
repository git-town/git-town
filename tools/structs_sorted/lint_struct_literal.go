package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"slices"
)

func lintStructLiteralVariable(compositeLit *ast.CompositeLit, fileSet *token.FileSet) Issues {
	structType, ok := compositeLit.Type.(*ast.Ident)
	if !ok {
		return Issues{}
	}
	structName := structType.Name
	if slices.Contains(ignoreTypes, structName) {
		return Issues{}
	}
	fieldNames := structInstantiationFieldNames(compositeLit)
	if len(fieldNames) == 0 {
		return Issues{}
	}
	sortedFields := make([]string, len(fieldNames))
	copy(sortedFields, fieldNames)
	slices.Sort(sortedFields)
	if reflect.DeepEqual(fieldNames, sortedFields) {
		return Issues{}
	}
	return Issues{
		issue{
			expected:   sortedFields,
			position:   fileSet.Position(compositeLit.Pos()),
			structName: structName,
		},
	}
}

func structInstantiationFieldNames(compositeLit *ast.CompositeLit) []string {
	result := make([]string, 0, len(compositeLit.Elts))
	for _, expr := range compositeLit.Elts {
		kvExpr, ok := expr.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kvExpr.Key.(*ast.Ident)
		if !ok {
			continue
		}
		result = append(result, ident.Name)
	}
	return result
}
