package main

import (
	"go/ast"
	"go/token"
	"slices"

	"github.com/git-town/git-town/tools/tests_sorted/matcher"
)

func lintFuncDecl(funcSpec *ast.FuncDecl, fileSet *token.FileSet, issues *Issues) error {
	return lintFunc(funcSpec.Pos(), funcSpec.Type, funcSpec.Body, fileSet, issues)
}

func lintFuncLit(funcSpec *ast.FuncLit, fileSet *token.FileSet, issues *Issues) error {
	return lintFunc(funcSpec.Pos(), funcSpec.Type, funcSpec.Body, fileSet, issues)
}

func lintFunc(pos token.Pos, funcType *ast.FuncType, funcBody *ast.BlockStmt, fileSet *token.FileSet, issues *Issues) error {
	if !isTestFunction(funcType) {
		return nil
	}

	subtests, err := topLevelRunNames(funcBody.List)
	if err != nil {
		return err
	}
	sortedSubtests := make([]string, len(subtests))
	copy(sortedSubtests, subtests)
	slices.Sort(sortedSubtests)

	if slices.Equal(subtests, sortedSubtests) {
		return nil
	}
	*issues = append(*issues, issue{
		expected: sortedSubtests,
		position: fileSet.Position(pos),
	})
	return nil
}

// isTestFunction returns true if a given funcSpec describes a test function/helper.
func isTestFunction(funcType *ast.FuncType) bool {
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

func topLevelRunNames(statements []ast.Stmt) ([]string, error) {
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
		subtest, r := testNameExtractor.Extract(expr.X)
		if !r.Success() {
			continue
		}
		subtests = append(subtests, subtest)
	}
	return subtests, nil
}
