package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"slices"
)

func lintStructLiteralCallArg(callExpr *ast.CallExpr, fileSet *token.FileSet) Issues {
	result := Issues{}
	for _, arg := range callExpr.Args {
		compositeLit, ok := arg.(*ast.CompositeLit)
		if !ok {
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
		result = append(result, issue{
			expected: sortedFields,
			position: fileSet.Position(compositeLit.Pos()),
		})
	}
	return result
}
