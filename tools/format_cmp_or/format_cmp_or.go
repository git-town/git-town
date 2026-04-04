package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	err := filepath.WalkDir(".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil || dirEntry.IsDir() || !isGoFile(path) || shouldIgnorePath(path) {
			return err
		}
		changed, err := formatFile(path, dirEntry.Type().Perm())
		if err != nil {
			return err
		}
		if changed {
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func formatFile(path string, perm fs.FileMode) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	newContent := FormatFileContent(path, content)
	if string(newContent) == string(content) {
		return false, nil
	}
	return true, os.WriteFile(path, newContent, perm)
}

// FormatFileContent sorts the arguments of all cmp.Or invocations in the given Go source.
// Calls that already contain keep-sorted markers are skipped since they are managed by the keep-sorted tool.
func FormatFileContent(path string, src []byte) []byte {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		fmt.Printf("Cannot parse file: %v\n", err)
		return src
	}
	var replacements []replacement
	ast.Inspect(file, func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if !isCmpOr(callExpr) {
			return true
		}
		replacements = append(replacements, computeReplacements(fset, file, callExpr)...)
		return true
	})
	if len(replacements) == 0 {
		return src
	}
	// Apply replacements from back to front to preserve earlier byte offsets.
	slices.SortFunc(replacements, func(a, b replacement) int {
		return b.start - a.start
	})
	result := make([]byte, len(src))
	copy(result, src)
	for _, replacement := range replacements {
		newResult := make([]byte, 0, len(result))
		newResult = append(newResult, result[:replacement.start]...)
		newResult = append(newResult, replacement.newText...)
		newResult = append(newResult, result[replacement.end:]...)
		result = newResult
	}
	return result
}

type replacement struct {
	end     int
	newText string
	start   int
}

func isCmpOr(callExpr *ast.CallExpr) bool {
	selector, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkg.Name == "cmp" && selector.Sel.Name == "Or"
}

func computeReplacements(fset *token.FileSet, file *ast.File, callExpr *ast.CallExpr) []replacement {
	if len(callExpr.Args) < 2 {
		return nil
	}
	// All arguments must be plain identifiers; skip the call otherwise.
	names := make([]string, len(callExpr.Args))
	for i, arg := range callExpr.Args {
		ident, ok := arg.(*ast.Ident)
		if !ok {
			return nil
		}
		names[i] = ident.Name
	}
	// report calls already managed by keep-sorted.
	if hasKeepSortedComment(fset, file, callExpr) {
		position := fset.Position(callExpr.Pos())
		fmt.Printf("%s:%d: remove keep-sorted marker from cmp.Or call\n", position.Filename, position.Line)
		return nil
	}
	sortedNames := make([]string, len(names))
	copy(sortedNames, names)
	slices.SortFunc(sortedNames, func(a, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})
	if slices.Equal(names, sortedNames) {
		return nil
	}
	var reps []replacement
	for i, arg := range callExpr.Args {
		if names[i] == sortedNames[i] {
			continue
		}
		reps = append(reps, replacement{
			end:     fset.Position(arg.End()).Offset,
			newText: sortedNames[i],
			start:   fset.Position(arg.Pos()).Offset,
		})
	}
	return reps
}

// hasKeepSortedComment reports whether the call expression contains a keep-sorted comment.
func hasKeepSortedComment(fset *token.FileSet, file *ast.File, callExpr *ast.CallExpr) bool {
	callStart := fset.Position(callExpr.Pos()).Offset
	callEnd := fset.Position(callExpr.End()).Offset
	for _, commentGroup := range file.Comments {
		offset := fset.Position(commentGroup.Pos()).Offset
		if offset < callStart || offset >= callEnd {
			continue
		}
		for _, comment := range commentGroup.List {
			if strings.Contains(comment.Text, "keep-sorted") {
				return true
			}
		}
	}
	return false
}

func isGoFile(path string) bool {
	return strings.HasSuffix(path, ".go")
}

func shouldIgnorePath(path string) bool {
	return strings.HasPrefix(path, "vendor/") ||
		strings.HasPrefix(path, ".git/")
}
