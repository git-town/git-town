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

// cmpEqualFinder implements the ast.Visitor interface to find calls to "cmp.Equal".
type cmpEqualFinder struct {
	filePath string         // path of the file being currently visited
	fileSet  *token.FileSet // position information for AST nodes in the current file
}

// Visit is called by ast.Walk for each node in the AST.
func (self *cmpEqualFinder) Visit(node ast.Node) ast.Visitor {
	// ensure the AST node is a function call expression
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return self
	}

	// ensure the function being called is a selector expression, e.g. "pkg.Func"
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return self
	}

	// ensure the package part of the selector is an identifier, e.g. "cmp"
	pkgIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return self
	}

	// ensure the package name is "cmp" and the function name is "Equal"
	if pkgIdent.Name == "cmp" && selectorExpr.Sel.Name == "Equal" {
		position := self.fileSet.Position(callExpr.Pos())
		fmt.Printf("%s:%d: Please call equal.Equal instead of cmp.Equal\n", self.filePath, position.Line)
	}

	return self
}

// parses the given Go file and walks its AST nodes to find cmp.Equal calls
func lintFile(filePath string) error {
	fileSet := token.NewFileSet() // holds position information
	fileAST, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file %s: %w", filePath, err)
	}
	visitor := &cmpEqualFinder{
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
	if strings.HasPrefix(cleanedPath, "pkg"+string(filepath.Separator)+"equal"+string(filepath.Separator)) {
		return true
	}
	return false
}

func main() {
	err := filepath.Walk(".", func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil || shouldSkipPath(path) || fileInfo.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
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
