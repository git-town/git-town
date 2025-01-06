package matcher

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/shoenig/test/must"
)

type trueMatcher struct{}

func (self *trueMatcher) Match(ast.Expr) Result {
	return okResult
}

type falseMatcher struct{}

func (self *falseMatcher) Match(ast.Expr) Result {
	return stringResult("fake failure")
}

type trueFieldMatcher struct{}

func (self *trueFieldMatcher) Match(PositionalField) Result {
	return okResult
}

type falseFieldMatcher struct{}

func (self *falseFieldMatcher) Match(PositionalField) Result {
	return stringResult("fake field failure")
}

func TestPositionalFields(t *testing.T) {
	expr, err := parser.ParseExpr(`func(a, b string, c bool)`)
	must.NoError(t, err)
	funcType := expr.(*ast.FuncType)

	var positionalIndices []int
	var positionalFields []PositionalField
	for i, field := range PositionalFields(funcType.Params) {
		positionalIndices = append(positionalIndices, i)
		positionalFields = append(positionalFields, field)
	}

	must.Eq(t, []int{0, 1, 2}, positionalIndices)
	must.Eq(t, len(positionalFields), 3)
	must.Eq(t, "a", positionalFields[0].Name.Name)
	must.Eq(t, "b", positionalFields[1].Name.Name)
	must.Eq(t, "c", positionalFields[2].Name.Name)
}

func TestIdentSelectorMatcher(t *testing.T) {
	t.Parallel()

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			expr       string
			matcher    *IdentSelectorMatcher
			wantReason string
		}{
			{
				expr:       `func()`,
				matcher:    &IdentSelectorMatcher{},
				wantReason: "not an ast.SelectorExpr",
			},
			{
				expr: `foo().bar`,
				matcher: &IdentSelectorMatcher{
					Namespace: "foo",
					Name:      "bar",
				},
				wantReason: "namespace is not an ast.Ident",
			},
			{
				expr: `notFoo.bar`,
				matcher: &IdentSelectorMatcher{
					Namespace: "foo",
					Name:      "bar",
				},
				wantReason: `namespace doesn't match: want "foo", got "notFoo"`,
			},
			{
				expr: `foo.notBar`,
				matcher: &IdentSelectorMatcher{
					Namespace: "foo",
					Name:      "bar",
				},
				wantReason: `name doesn't match: want "bar", got "notBar"`,
			},
		} {
			t.Run(tc.wantReason, func(t *testing.T) {
				expr, err := parser.ParseExpr(tc.expr)
				must.NoError(t, err)

				r := tc.matcher.Match(expr)

				must.False(t, r.Success())
				must.Eq(t, tc.wantReason, r.FailureReason())
			})
		}
		t.Run("Success", func(t *testing.T) {
			expr, err := parser.ParseExpr(`foo.bar`)
			must.NoError(t, err)
			m := &IdentSelectorMatcher{
				Namespace: "foo",
				Name:      "bar",
			}

			r := m.Match(expr)

			must.True(t, r.Success())
		})
	})
}

func TestPointerMatcher(t *testing.T) {
	t.Parallel()

	t.Run("InnerMatchFailure", func(t *testing.T) {
		expr, err := parser.ParseExpr(`*foo`)
		must.NoError(t, err)
		m := &PointerMatcher{
			InnerMatcher: &falseMatcher{},
		}

		r := m.Match(expr)

		must.False(t, r.Success())
		must.Eq(t, "fake failure", r.FailureReason())
	})

	t.Run("NotStarExpr", func(t *testing.T) {
		expr, err := parser.ParseExpr(`Foo`)
		must.NoError(t, err)
		m := &PointerMatcher{
			InnerMatcher: &trueMatcher{},
		}

		r := m.Match(expr)

		must.False(t, r.Success())
		must.Eq(t, "not an ast.StarExpr", r.FailureReason())
	})

	t.Run("Success", func(t *testing.T) {
		expr, err := parser.ParseExpr(`*foo`)
		must.NoError(t, err)
		m := &PointerMatcher{
			InnerMatcher: &trueMatcher{},
		}

		r := m.Match(expr)

		must.True(t, r.Success())
	})
}

func firstFuncParam(t *testing.T, funcTypeText string) PositionalField {
	t.Helper()
	expr, err := parser.ParseExpr(funcTypeText)
	must.NoError(t, err)
	funcType := expr.(*ast.FuncType)
	for _, field := range PositionalFields(funcType.Params) {
		return field
	}
	must.Unreachable(t)
	return PositionalField{}
}

func TestFieldMatcher(t *testing.T) {
	t.Parallel()

	t.Run("Failure", func(t *testing.T) {
		field := firstFuncParam(t, `func(notA string)`)
		m := &FieldMatcher{
			Name:        "a",
			TypeMatcher: &trueMatcher{},
		}

		r := m.Match(field)

		must.False(t, r.Success())
		must.Eq(t, `field name doesn't match: want "a", got "notA"`, r.FailureReason())
	})

	t.Run("Success", func(t *testing.T) {
		field := firstFuncParam(t, `func(a string)`)
		m := &FieldMatcher{
			Name:        "a",
			TypeMatcher: &trueMatcher{},
		}

		r := m.Match(field)

		must.True(t, r.Success())
	})
}

func funcParams(t *testing.T, funcTypeText string) *ast.FieldList {
	t.Helper()
	expr, err := parser.ParseExpr(funcTypeText)
	must.NoError(t, err)
	funcType := expr.(*ast.FuncType)
	return funcType.Params
}

func TestFieldListPrefixMatcher(t *testing.T) {
	t.Parallel()

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()
		t.Run("NotEnoughFields", func(t *testing.T) {
			params := funcParams(t, `func()`)
			m := &FieldListPrefixMatcher{
				Prefix: []PositionalFieldMatcher{
					&trueFieldMatcher{},
					&trueFieldMatcher{},
				},
			}

			r := m.Match(params)

			must.False(t, r.Success())
			must.Eq(t, "not enough fields, want at least 2, got 0", r.FailureReason())
		})
		t.Run("FieldDoesNotMatch", func(t *testing.T) {
			params := funcParams(t, `func(a string)`)
			m := &FieldListPrefixMatcher{
				Prefix: []PositionalFieldMatcher{
					&falseFieldMatcher{},
				},
			}

			r := m.Match(params)

			must.False(t, r.Success())
			must.Eq(t, "field[0] doesn't match: fake field failure", r.FailureReason())
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			desc    string
			params  *ast.FieldList
			matcher *FieldListPrefixMatcher
		}{
			{
				desc:   "EqualLenMatchers",
				params: funcParams(t, `func(a, b, c string)`),
				matcher: &FieldListPrefixMatcher{
					Prefix: []PositionalFieldMatcher{
						&trueFieldMatcher{},
						&trueFieldMatcher{},
						&trueFieldMatcher{},
					},
				},
			}, {
				desc:   "NoMatchersNoParams",
				params: funcParams(t, `func()`),
				matcher: &FieldListPrefixMatcher{
					Prefix: []PositionalFieldMatcher{},
				},
			}, {
				desc:   "SomeMatchers",
				params: funcParams(t, `func(a, b, c string)`),
				matcher: &FieldListPrefixMatcher{
					Prefix: []PositionalFieldMatcher{
						&trueFieldMatcher{},
						&trueFieldMatcher{},
					},
				},
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				r := tc.matcher.Match(tc.params)

				must.True(t, r.Success())
			})
		}
	})
}

func TestFirstStringArgFromFuncCallExtractor(t *testing.T) {
	t.Parallel()

	t.Run("Failure", func(t *testing.T) {
		for _, tc := range []struct {
			funcMatcher ExprMatcher
			expr        string
			wantReason  string
		}{
			{
				funcMatcher: &trueMatcher{},
				expr:        `struct{}`,
				wantReason:  "not an ast.CallExpr",
			}, {
				funcMatcher: &falseMatcher{},
				expr:        `foo("arg1", "arg2")`,
				wantReason:  "call.Fun doesn't match: fake failure",
			}, {
				funcMatcher: &trueMatcher{},
				expr:        `foo()`,
				wantReason:  "len(call.Args) == 0, want at least 1",
			}, {
				funcMatcher: &trueMatcher{},
				expr:        `foo(""+"")`,
				wantReason:  "the first call argument is not an ast.BasicLit",
			}, {
				funcMatcher: &trueMatcher{},
				expr:        `foo(1)`,
				wantReason:  "the first call argument is not a STRING ast.BasicLit",
			},
		} {
			t.Run(tc.wantReason, func(t *testing.T) {
				expr, err := parser.ParseExpr(tc.expr)
				must.NoError(t, err)
				e := &FirstStringArgFromFuncCallExtractor{
					FuncMatcher: tc.funcMatcher,
				}

				_, r, err := e.Extract(expr)

				must.NoError(t, err)
				must.False(t, r.Success())
				must.Eq(t, tc.wantReason, r.FailureReason())
			})
		}
		t.Run("UnquoteError", func(t *testing.T) {
			// Construct a correct string literal, and then make it invalid by assigning
			// a known badly quoted string to its Value.
			expr, err := parser.ParseExpr(`foo("test")`)
			must.NoError(t, err)
			call := expr.(*ast.CallExpr)
			arg := call.Args[0].(*ast.BasicLit)
			arg.Value = `"foo` // Intentionally badly quoted for this test.
			e := &FirstStringArgFromFuncCallExtractor{
				FuncMatcher: &trueMatcher{},
			}

			_, _, err = e.Extract(expr)

			must.ErrorContains(t, err, "the first call argument is an invalid string literal:")
		})
	})

	t.Run("Success", func(t *testing.T) {
		expr, err := parser.ParseExpr(`foo("arg1")`)
		must.NoError(t, err)
		e := &FirstStringArgFromFuncCallExtractor{
			FuncMatcher: &trueMatcher{},
		}

		arg, r, err := e.Extract(expr)

		must.NoError(t, err)
		must.True(t, r.Success())
		must.Eq(t, "arg1", arg)
	})
}
