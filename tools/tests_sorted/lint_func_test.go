package main

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/shoenig/test/must"
)

func funcType(t *testing.T, funcText string) *ast.FuncType {
	expr, err := parser.ParseExpr(funcText)
	must.NoError(t, err)
	return expr.(*ast.FuncType)
}

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
			fType := funcType(t, tc.expr)

			got := isTestFunc(fType)

			must.Eq(t, tc.want, got)
		})
	}
}

func funcStatements(t *testing.T, funcText string) []ast.Stmt {
	expr, err := parser.ParseExpr(funcText)
	must.NoError(t, err)
	funcLit := expr.(*ast.FuncLit)
	return funcLit.Body.List
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
			statements := funcStatements(t, tc.expr)

			got, err := tRunSubtestNames(statements)

			must.NoError(t, err)
			must.Eq(t, tc.wantSubtestNames, got)
		})
	}
}
