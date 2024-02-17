package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"slices"
)

func lintStructLiteralCallArg(callExpr *ast.CallExpr, fileSet *token.FileSet, issues *Issues) {
	for _, arg := range callExpr.Args {
		compositeLit, ok := arg.(*ast.CompositeLit)
		if !ok {
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
}
