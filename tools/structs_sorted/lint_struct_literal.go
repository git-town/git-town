package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"slices"
)

func lintStructLiteralVariable(compositeLit *ast.CompositeLit, fileSet *token.FileSet, issues *Issues) {
	structType, ok := compositeLit.Type.(*ast.Ident)
	if !ok {
		return
	}
	structName := structType.Name
	if slices.Contains(ignoreTypes, structName) {
		return
	}
	fieldNames := structInstantiationFieldNames(compositeLit)
	if len(fieldNames) == 0 {
		return
	}
	sortedFields := make([]string, len(fieldNames))
	copy(sortedFields, fieldNames)
	slices.Sort(sortedFields)
	if reflect.DeepEqual(fieldNames, sortedFields) {
		return
	}
	*issues = append(*issues, issue{
		expected: sortedFields,
		position: fileSet.Position(compositeLit.Pos()),
	})
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
