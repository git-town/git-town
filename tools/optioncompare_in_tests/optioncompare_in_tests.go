package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// wrongCompareFinder implements the ast.Visitor interface to find and transform calls to "must.Eq(t, Some(x), y)".
type wrongCompareFinder struct {
	filePath string         // path of the file being currently visited
	fileSet  *token.FileSet // position information for AST nodes in the current file
	modified bool           // tracks whether any modifications were made
}

// Visit is called by ast.Walk for each node in the AST.
func (self *wrongCompareFinder) Visit(node ast.Node) ast.Visitor {
	// ensure the AST node is a function call expression
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return self
	}

	// ensure the function being called is a selector expression, e.g. "must.Eq"
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return self
	}

	// ensure the package part of the selector is an identifier, e.g. "must"
	pkgIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return self
	}

	// ensure the package name is "must" and the function name is "Eq"
	if pkgIdent.Name != "must" || selectorExpr.Sel.Name != "Eq" {
		return self
	}

	// ensure if we have at least 3 arguments
	if len(callExpr.Args) < 3 {
		return self
	}

	// ensure if the second argument is a call to "Some"
	someCallExpr := getSomeCall(callExpr.Args[1])
	if someCallExpr == nil {
		return self
	}

	// Transform must.Eq(t, Some(x), y) to must.True(t, y.EqualSome(x))
	transformToEqualSome(callExpr, someCallExpr, callExpr.Args[2])
	self.modified = true

	return self
}

// getSomeCall checks if the given expression is a call to "Some" or "Some[T]" and returns the call expression if so
func getSomeCall(expr ast.Expr) *ast.CallExpr {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil
	}

	// check if the function being called is an identifier "Some"
	ident, ok := callExpr.Fun.(*ast.Ident)
	if ok && ident.Name == "Some" {
		return callExpr
	}

	// check if the function being called is a generic type "Some[T]"
	if indexExpr, ok := callExpr.Fun.(*ast.IndexExpr); ok {
		ident, ok := indexExpr.X.(*ast.Ident)
		if ok && ident.Name == "Some" {
			return callExpr
		}
	}

	return nil
}

// isTestFile checks if the file is a unit test file
func isTestFile(filePath string) bool {
	return strings.HasSuffix(filePath, "_test.go") || filePath == "internal/test/cucumber/steps.go"
}

// parses the given Go file and walks its AST nodes to transform the pattern
func lintFile(filePath string) error {
	fileSet := token.NewFileSet() // holds position information
	fileAST, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file %s: %w", filePath, err)
	}
	visitor := &wrongCompareFinder{
		filePath: filePath,
		fileSet:  fileSet,
		modified: false,
	}
	ast.Walk(visitor, fileAST)

	// If modifications were made, write the modified AST back to the file
	if visitor.modified {
		var buf bytes.Buffer
		if err := format.Node(&buf, fileSet, fileAST); err != nil {
			return fmt.Errorf("error formatting modified file %s: %w", filePath, err)
		}

		if err := os.WriteFile(filePath, buf.Bytes(), 0o600); err != nil {
			return fmt.Errorf("error writing modified file %s: %w", filePath, err)
		}

		fmt.Printf("%s: fixed must.Eq(t, Some(x), y) patterns\n", filePath)
	}

	return nil
}

// indicates whether a given path should be skipped
func shouldSkipPath(path string) bool {
	cleanedPath := filepath.Clean(path) // resolve any ".." or "." components and get a canonical path.
	if strings.HasPrefix(cleanedPath, "vendor"+string(filepath.Separator)) {
		return true
	}
	return false
}

// transformToEqualSome transforms must.Eq(t, Some(x), y) to must.True(t, y.EqualSome(x))
func transformToEqualSome(mustEqCall *ast.CallExpr, someCall *ast.CallExpr, yArg ast.Expr) {
	// Change the selector from "Eq" to "True"
	selectorExpr := mustEqCall.Fun.(*ast.SelectorExpr)
	selectorExpr.Sel.Name = "True"

	// Create y.EqualSome(x) where x is the argument to Some
	equalSomeCall := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: yArg,
			Sel: &ast.Ident{
				Name: "EqualSome",
			},
		},
		Args: someCall.Args, // Use the arguments from the Some() call
	}

	// Replace the arguments: keep t (first arg), replace second and third with y.EqualSome(x)
	mustEqCall.Args = []ast.Expr{
		mustEqCall.Args[0], // keep t
		equalSomeCall,      // y.EqualSome(x)
	}
}

func main() {
	err := filepath.Walk(".", func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil || shouldSkipPath(path) || fileInfo.IsDir() || !strings.HasSuffix(path, ".go") || !isTestFile(path) {
			return err
		}
		if err := lintFile(path); err != nil {
			return fmt.Errorf("error linting file %s: %w", path, err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking current directory: %v\n", err)
	}
}
