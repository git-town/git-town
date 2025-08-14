package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// somePatternFinder implements the ast.Visitor interface to find calls to "must.Eq(t, Some(x), y)".
type somePatternFinder struct {
	filePath string         // path of the file being currently visited
	fileSet  *token.FileSet // position information for AST nodes in the current file
}

// Visit is called by ast.Walk for each node in the AST.
func (self *somePatternFinder) Visit(node ast.Node) ast.Visitor {
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
	if pkgIdent.Name == "must" && selectorExpr.Sel.Name == "Eq" {
		// check if we have at least 3 arguments
		if len(callExpr.Args) >= 3 {
			// check if the second argument is a call to "Some"
			if self.isSomeCall(callExpr.Args[1]) {
				position := self.fileSet.Position(callExpr.Pos())
				fmt.Printf("%s:%d: must.Eq(t, Some(x), y) pattern detected\n", self.filePath, position.Line)
			}
		}
	}

	return self
}

// isSomeCall checks if the given expression is a call to "Some"
func (self *somePatternFinder) isSomeCall(expr ast.Expr) bool {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	// check if the function being called is an identifier "Some"
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "Some"
}

// isTestFile checks if the file is a unit test file
func isTestFile(filePath string) bool {
	return strings.HasSuffix(filePath, "_test.go")
}

// parses the given Go file and walks its AST nodes to find the pattern
func lintFile(filePath string) error {
	fileSet := token.NewFileSet() // holds position information
	fileAST, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file %s: %w", filePath, err)
	}
	visitor := &somePatternFinder{
		filePath: filePath,
		fileSet:  fileSet,
	}
	ast.Walk(visitor, fileAST)
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

func main() {
	err := filepath.Walk(".", func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil || shouldSkipPath(path) || fileInfo.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		// only lint test files
		if !isTestFile(path) {
			return nil
		}
		if err := lintFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error linting file %s: %v\n", path, err)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking current directory: %v\n", err)
	}
}
