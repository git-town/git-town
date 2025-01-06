package main

import (
	"go/ast"

	"github.com/git-town/git-town/tools/tests_sorted/matcher"
)

// isTestFunc returns true if a given funcType describes a test function/helper.
//
// `ast.FuncType` does not differentiate between top-level functions, function
// literals, or between tests or test helpers. This is merely a signature check.
// The main objective of isTestFunc is to select test functions that are
// compatible with `tRunSubtestNames` as we make an assumption about the first
// parameter being `t *testing.T` named `t` specifically.
func isTestFunc(funcType *ast.FuncType) bool {
	firstParam := &matcher.FieldMatcher{
		Name: "t",
		TypeMatcher: &matcher.PointerMatcher{
			InnerMatcher: &matcher.IdentSelectorMatcher{
				Namespace: "testing",
				Name:      "T",
			},
		},
	}
	m := &matcher.FieldListPrefixMatcher{
		Prefix: []matcher.PositionalFieldMatcher{firstParam},
	}
	return m.Match(funcType.Params).Success()
}

// tRunSubtestNames returns a list of subtest names from `t.Run(...)`
// invocations.
func tRunSubtestNames(statements []ast.Stmt) ([]string, error) {
	var subtests []string
	testNameExtractor := &matcher.FirstStringArgFromFuncCallExtractor{
		FuncMatcher: &matcher.IdentSelectorMatcher{
			Namespace: "t",
			Name:      "Run",
		},
	}
	for _, stmt := range statements {
		expr, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		subtest, r, err := testNameExtractor.Extract(expr.X)
		if err != nil {
			// Failure to extract should be reported.
			// This may indicate a linter bug.
			return nil, err
		}
		if !r.Success() {
			continue
		}
		subtests = append(subtests, subtest)
	}
	return subtests, nil
}
