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

// cmpEqualVisitor implements the ast.Visitor interface to find calls to "cmp.Equal".
type cmpEqualVisitor struct {
	fileSet  *token.FileSet // FileSet to get position information for nodes.
	filePath string         // Path of the file being currently visited.
}

// Visit is called for each node in the AST.
func (v *cmpEqualVisitor) Visit(node ast.Node) ast.Visitor {
	// Check if the node is a function call expression.
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		// If it's not a call expression, continue traversing its children.
		return v
	}

	// Check if the function being called is a selector expression (e.g., "pkg.Func").
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		// If it's not a selector expression, continue.
		return v
	}

	// Check if the package part of the selector is an identifier (e.g., "cmp").
	pkgIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		// If the package part is not an identifier, continue.
		return v
	}

	// Check if the package name is "cmp" and the function name is "Equal".
	if pkgIdent.Name == "cmp" && selectorExpr.Sel.Name == "Equal" {
		// Get the position (line number) of the call.
		position := v.fileSet.Position(callExpr.Pos())
		fmt.Printf("%s:%d: Please call equal.Equal instead of cmp.Equal\n", v.filePath, position.Line)
	}

	// Continue visiting children of the current node.
	return v
}

// lintFile parses a single Go file and walks its AST to find cmp.Equal calls.
func lintFile(filePath string) error {
	// Create a new FileSet to hold position information.
	fset := token.NewFileSet()

	// Parse the Go source file.
	// parser.ParseFile returns the AST of the file.
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file %s: %w", filePath, err)
	}

	// Create a new visitor instance for this file.
	visitor := &cmpEqualVisitor{
		fileSet:  fset,
		filePath: filePath,
	}

	// Walk the AST, applying our visitor to each node.
	ast.Walk(visitor, node)

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
	// Determine the paths to lint.
	// If no arguments are provided, lint the current directory recursively.
	// Otherwise, lint the provided arguments (files or directories).
	var pathsToLint []string
	if len(os.Args) > 1 {
		pathsToLint = os.Args[1:]
	} else {
		pathsToLint = []string{"."} // Default to current directory
	}

	// Iterate over the provided paths.
	for _, p := range pathsToLint {
		info, err := os.Stat(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %s: %v\n", p, err)
			continue
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
				fmt.Fprintf(os.Stderr, "Skipping file: %s\n", p)
				continue
			}
			if err := lintFile(p); err != nil {
				fmt.Fprintf(os.Stderr, "Error linting file %s: %v\n", p, err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Skipping non-.go file or unsupported path: %s\n", p)
		}
	}
}
