package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// cmpEqualVisitor implements the ast.Visitor interface to find calls to "cmp.Equal".
type cmpEqualVisitor struct {
	fileSet  *token.FileSet // FileSet to get position information for nodes.
	filePath string         // Path of the file being currently visited.
}

// Visit is called for each node in the AST.
func (v *cmpEqualVisitor) Visit(node ast.Node) ast.Visitor {
	// ensure the node is a function call expression
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return v
	}

	// ensure the function being called is a selector expression (e.g., "pkg.Func")
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	// ensure the package part of the selector is an identifier (e.g., "cmp")
	pkgIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return v
	}

	// ensure the package name is "cmp" and the function name is "Equal"
	if pkgIdent.Name == "cmp" && selectorExpr.Sel.Name == "Equal" {
		position := v.fileSet.Position(callExpr.Pos())
		fmt.Printf("%s:%d: Please call equal.Equal instead of cmp.Equal\n", v.filePath, position.Line)
	}

	return v
}

// lintFile parses a single Go file and walks all its AST nodes to find cmp.Equal calls.
func lintFile(filePath string) error {
	fileSet := token.NewFileSet() // holds position information
	fileAST, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file %s: %w", filePath, err)
	}
	visitor := &cmpEqualVisitor{
		fileSet:  fileSet,
		filePath: filePath,
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
	p := "."
	info, err := os.Stat(p)
	if err != nil {
		log.Fatalf("Error accessing path %s: %v\n", p, err)
	}

	if info.IsDir() {
		// If it's a directory, walk it recursively.
		err := filepath.Walk(p, func(path string, d os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if the current path (file or directory) should be skipped.
			if shouldSkipPath(path) {
				return nil // Return filepath.SkipDir if it's a vendor or pkg/equal directory, or nil to skip file.
			}

			// Only process .go files and skip directories themselves.
			if !d.IsDir() && strings.HasSuffix(path, ".go") {
				if err := lintFile(path); err != nil {
					fmt.Fprintf(os.Stderr, "Error linting file %s: %v\n", path, err)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", p, err)
		}
	} else if strings.HasSuffix(p, ".go") {
		// If it's a .go file, lint it directly.
		// Check if the file path itself is within a vendor or pkg/equal directory.
		if shouldSkipPath(p) { // We know it's not a directory here
			log.Fatalf("Skipping file: %s\n", p)
		}
		if err := lintFile(p); err != nil {
			fmt.Fprintf(os.Stderr, "Error linting file %s: %v\n", p, err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Skipping non-.go file or unsupported path: %s\n", p)
	}
}
