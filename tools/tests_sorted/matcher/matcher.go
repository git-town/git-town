// Package matcher declarative go/ast matchers.
package matcher

import (
	"fmt"
	"go/ast"
	"go/token"
	"iter"
	"strconv"
)

// Result is a match result with a reason.
// Use Result.Success() to convert it to bool.
// FailureReason() is useful for unit tests.
type Result interface {
	Success() bool
	FailureReason() string
}

type stringResult string

func (self stringResult) Success() bool {
	return self == ""
}

func (self stringResult) FailureReason() string {
	return string(self)
}

const okResult = stringResult("")

func fmtResult(format string, a ...any) stringResult {
	return stringResult(fmt.Sprintf(format, a...))
}

// ExprMatcher matches ast.Expr.
type ExprMatcher interface {
	Match(expr ast.Expr) Result
}

// PositionalField represents a logical field from an ast.FieldList.
// Go allows fields and function parameters to be grouped by type, e.g.
// `func(a, b string)`, which comes up as a single `ast.Field` element with two
// `ast.Field.Names`. PositionalField represents a single logical `a string`
// regardless of syntactic grouping.
type PositionalField struct {
	Name  *ast.Ident
	Field *ast.Field
}

// PositionalFieldMatcher matches a PositionalField.
type PositionalFieldMatcher interface {
	Match(field PositionalField) Result
}

// FieldListMatcher matches an ast.FieldList.
type FieldListMatcher interface {
	Match(fields *ast.FieldList) Result
}

// PositionalFields yields (name, field) for each name in ast.FieldList in logical order.
// For example func(a, b string, c int) will result in
// - fields[0]: ast.Field{Names: []{"a", "b"}}
// - fields[1]: ast.Field{Names: []{"c"}}
// So PositionalFields(fields) will yield
// - ("a", fields[0])
// - ("b", fields[0])
// - ("c", fields[1])
func PositionalFields(fields *ast.FieldList) iter.Seq2[int, PositionalField] {
	return func(yield func(int, PositionalField) bool) {
		i := 0
		for _, field := range fields.List {
			for _, name := range field.Names {
				if !yield(i, PositionalField{Name: name, Field: field}) {
					return // Stopping early.
				}
				i++
			}
		}
	}
}

// IdentSelectorMatcher matches ast.SelectorExpr where both the selector
// expression (left of the `.`) and the selector itself (right of the `.`) are
// given `Namespace.Name`.
type IdentSelectorMatcher struct {
	Namespace string
	Name      string
}

func (self *IdentSelectorMatcher) Match(expr ast.Expr) Result {
	// Check that the whole expression is `<namespace>.<name>`.
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return stringResult("not an ast.SelectorExpr")
	}
	// Check that the `<namespace>` is an identifier.
	namespace, ok := selector.X.(*ast.Ident)
	if !ok {
		return stringResult("namespace is not an ast.Ident")
	}
	// Match the names.
	if self.Namespace != namespace.Name {
		return fmtResult("namespace doesn't match: want %q, got %q", self.Namespace, namespace.Name)
	}
	if self.Name != selector.Sel.Name {
		return fmtResult("name doesn't match: want %q, got %q", self.Name, selector.Sel.Name)
	}
	return okResult
}

// PointerMatcher matches an ast.StarExpr where the inner expression matcher the InnerMatcher.
type PointerMatcher struct {
	InnerMatcher ExprMatcher
}

func (self *PointerMatcher) Match(expr ast.Expr) Result {
	// Verify that this is a `*<expr>`.
	ptr, ok := expr.(*ast.StarExpr)
	if !ok {
		return stringResult("not an ast.StarExpr")
	}
	// Delegate to the InnerMatcher.
	return self.InnerMatcher.Match(ptr.X)
}

// FieldMatcher matches an ast.Field element if it has a given Name and Type.
type FieldMatcher struct {
	Name        string
	TypeMatcher ExprMatcher
}

func (self *FieldMatcher) Match(field PositionalField) Result {
	if self.Name != field.Name.Name {
		return fmtResult("field name doesn't match: want %q, got %q", self.Name, field.Name.Name)
	}
	if r := self.TypeMatcher.Match(field.Field.Type); !r.Success() {
		return fmtResult("field type doesn't match: %v", r)
	}
	return okResult
}

// FieldListPrefixMatcher matches an `ast.FieldList` if
// - a field list contains at least len(Prefix) elements,
// - each ast.FieldList logical element i matches Prefix[i].
type FieldListPrefixMatcher struct {
	Prefix []PositionalFieldMatcher
}

func (self *FieldListPrefixMatcher) Match(fields *ast.FieldList) Result {
	var posFields []PositionalField
	for _, field := range PositionalFields(fields) {
		posFields = append(posFields, field)
	}
	for i, fieldMatcher := range self.Prefix {
		if i >= len(posFields) {
			return fmtResult("not enough fields, want at least %d, got %d", len(self.Prefix), len(posFields))
		}
		field := posFields[i]
		if r := fieldMatcher.Match(field); !r.Success() {
			return fmtResult("field[%d] doesn't match: %v", i, r.FailureReason())
		}
	}
	return okResult
}

// FirstStringArgFromFuncCallExtractor is an ast.CallExpr <Function(Arg1, ...)>
// extractor that extracts the first argument if the function matches the
// FuncMatcher and the first argument is a string literal.
//
// For example for `t.Run("foo", ...)` Extract() returns "foo", without quotes.
type FirstStringArgFromFuncCallExtractor struct {
	FuncMatcher ExprMatcher
}

func (self *FirstStringArgFromFuncCallExtractor) Extract(expr ast.Expr) (string, Result, error) {
	// Check that this is a <Fun>(<Call>) expression.
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return "", stringResult("not an ast.CallExpr"), nil
	}
	// Check that the call.Fun matches the FuncMatcher.
	if r := self.FuncMatcher.Match(call.Fun); !r.Success() {
		return "", fmtResult("call.Fun doesn't match: %v", r.FailureReason()), nil
	}
	// Check call.Args.
	if len(call.Args) < 1 {
		return "", fmtResult("len(call.Args) == %d, want at least 1", len(call.Args)), nil
	}
	firstArg := call.Args[0]
	// Check the first arg type.
	lit, ok := firstArg.(*ast.BasicLit)
	if !ok {
		return "", stringResult("the first call argument is not an ast.BasicLit"), nil
	}
	if lit.Kind != token.STRING {
		return "", stringResult("the first call argument is not a STRING ast.BasicLit"), nil
	}
	// The literal as is in code, including the quotes.
	quotedLiteral := lit.Value
	// https://go-review.googlesource.com/c/go/+/244960
	literal, err := strconv.Unquote(quotedLiteral)
	if err != nil {
		return "", okResult, fmt.Errorf("the first call argument is an invalid string literal: %v", err)
	}
	return literal, okResult, nil
}
