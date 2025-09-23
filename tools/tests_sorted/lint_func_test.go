package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/git-town/git-town/v22/pkg/equal"
	"github.com/google/go-cmp/cmp"
	"github.com/shoenig/test/must"
)

func TestIsTestFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		desc string
		expr string
		want bool
	}{
		{
			desc: "MultipleParameters",
			expr: `func(t *testing.T, x bool)`,
			want: true,
		}, {
			desc: "MultipleTParameters",
			expr: `func(t *testing.T, x *testing.T)`,
			want: true,
		}, {
			desc: "NotT",
			expr: `func(f *testing.T)`,
			want: false,
		}, {
			desc: "NotTestingT",
			expr: `func(t *string)`,
			want: false,
		}, {
			desc: "SingleParameter",
			expr: `func(t *testing.T)`,
			want: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			fType := funcType(t, tc.expr)

			got := isTestFunc(fType)

			must.Eq(t, tc.want, got)
		})
	}
}

func TestLintFuncDecl(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		desc       string
		funcText   string
		wantIssues Issues
	}{
		{
			desc: "IgnoresOtherCalls",
			funcText: `
				func TestFoo(t *testing.T) {
					t.Run("Bar")
					foo.Run("Another")
					t.NotRun("Unsorted")
					t.Run("Sum"+"Call")
					t.Run("Foo")
				}`,
			wantIssues: nil,
		},
		{
			desc: "RightOrder",
			funcText: `
				func TestFoo(t *testing.T) {
					t.Run("Bar")
					t.Run("Foo")
				}`,
			wantIssues: nil,
		},
		{
			desc: "WrongOrder",
			funcText: `
				func TestFoo(t *testing.T) {
					t.Run("Foo")
					t.Run("Bar")
				}`,
			wantIssues: Issues{
				{
					expected: []string{
						"Bar",
						"Foo",
					},
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			fset := token.NewFileSet()
			funcDecl := parseFuncDecl(t, fset, tc.funcText)

			issues, err := lintFuncDecl(funcDecl, fset)

			must.NoError(t, err)
			must.Eq(t, tc.wantIssues, issues, ignoreIssuePosition())
			// Defence in-depth. Verify that ignoreIssuePosition doesn't somehow entirely
			// break the comparison. must.SliceLen check must be after must.Eq to
			// preserve higher quality error messages from must.Eq.
			must.SliceLen(t, len(tc.wantIssues), issues)
		})
	}
}

func TestTRunSubtestNames(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		desc             string
		expr             string
		wantSubtestNames []string
	}{
		{
			desc: "Multiple",
			expr: `
				func() {
					t.Run("foo")
					t.Run("bar")
				}`,
			wantSubtestNames: []string{"foo", "bar"},
		},
		{
			desc: "Single t.Run",
			expr: `
				func() {
					t.Run("foo")
				}`,
			wantSubtestNames: []string{"foo"},
		},
		{
			desc: "Skips non-t.Run",
			expr: `
				func() {
					foo.Run("foo")
				}`,
			wantSubtestNames: nil,
		},
		{
			desc: "Skips t.not-Run",
			expr: `
				func() {
					t.Foo("foo")
				}`,
			wantSubtestNames: nil,
		},
		{
			desc: "Skips non-literals",
			expr: `
				func() {
					t.Run(""+"")
				}`,
			wantSubtestNames: nil,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			statements := funcStatements(t, tc.expr)

			got, err := tRunSubtestNames(statements)

			must.NoError(t, err)
			must.Eq(t, tc.wantSubtestNames, got)
		})
	}
}

func funcStatements(t *testing.T, funcText string) []ast.Stmt {
	t.Helper()
	expr, err := parser.ParseExpr(funcText)
	must.NoError(t, err)
	funcLit := expr.(*ast.FuncLit)
	return funcLit.Body.List
}

func funcType(t *testing.T, funcText string) *ast.FuncType {
	t.Helper()
	expr, err := parser.ParseExpr(funcText)
	must.NoError(t, err)
	return expr.(*ast.FuncType)
}

// ignoreIssuePosition returns a must.Setting that causes issues to be compared
// without taking issue.position into account.
func ignoreIssuePosition() must.Setting {
	return must.Cmp(cmp.Comparer(func(a, b issue) bool {
		return equal.Equal(a.expected, b.expected)
	}))
}

func parseFuncDecl(t *testing.T, fset *token.FileSet, funcText string) *ast.FuncDecl {
	t.Helper()
	// Package declaration is mandatory.
	// Concatenate the package and the function declaration with ';' to preserve
	// line numbers in any compilation errors.
	funcFile := "package test; " + funcText
	file, err := parser.ParseFile(fset, "test.go", funcFile, 0)
	must.NoError(t, err)
	must.SliceLen(t, 1, file.Decls) // We expect a single function.
	return file.Decls[0].(*ast.FuncDecl)
}
