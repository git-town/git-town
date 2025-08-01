package main

import (
	"go/ast"
	"go/token"
	"slices"

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

// lintFuncDecl returns a list of issues in a given ast.FuncDecl.
func lintFuncDecl(funcDecl *ast.FuncDecl, fileSet *token.FileSet) (Issues, error) {
	return lintFuncImpl(fileSet.Position(funcDecl.Pos()), funcDecl.Type, funcDecl.Body)
}

// lintFuncImpl returns a list of issues in a function described by its
// position, type and body.
func lintFuncImpl(pos token.Position, funcType *ast.FuncType, funcBody *ast.BlockStmt) (Issues, error) {
	if !isTestFunc(funcType) {
		return nil, nil
	}

	subtests, err := tRunSubtestNames(funcBody.List)
	if err != nil {
		return nil, err
	}
	sortedSubtests := make([]string, len(subtests))
	copy(sortedSubtests, subtests)
	slices.Sort(sortedSubtests)

	if !slices.Equal(subtests, sortedSubtests) {
		return Issues{
			issue{
				expected: sortedSubtests,
				position: pos,
			},
		}, nil
	}
	return nil, nil
}

// tRunSubtestNames returns a list of subtest names from `t.Run(...)`
// invocations.
func tRunSubtestNames(statements []ast.Stmt) ([]string, error) {
	var subtests []string //nolint:prealloc // alexkohler/prealloc#16
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
