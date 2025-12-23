package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/git-town/git-town/v22/pkg/asserts"
	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{
		Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs := asserts.NoError1(packages.Load(cfg, "./..."))
	foundError := false
	for _, pkg := range pkgs {
		for i, file := range pkg.Syntax {
			if shouldIgnorePath(pkg.GoFiles[i]) {
				continue
			}
			visitor := &addfVisitor{
				foundError: &foundError,
				fileSet:    pkg.Fset,
				path:       pkg.GoFiles[i],
				typeInfo:   pkg.TypesInfo,
			}
			ast.Walk(visitor, file)
		}
	}
	if foundError {
		os.Exit(1)
	}
}

type addfVisitor struct {
	foundError *bool
	fileSet    *token.FileSet
	path       string
	typeInfo   *types.Info
}

func (self *addfVisitor) Visit(node ast.Node) ast.Visitor {
	callExpr, isCallExpr := node.(*ast.CallExpr)
	if !isCallExpr {
		return self
	}

	// ensure this is a method call
	selectorExpr, isSelectorExpr := callExpr.Fun.(*ast.SelectorExpr)
	if !isSelectorExpr {
		return self
	}

	// ensure the receiver is of type stringslice.Collector
	if !self.isCollectorType(selectorExpr.X) {
		return self
	}

	switch selectorExpr.Sel.Name {
	case "Add":
		self.verifyAddCall(callExpr)
	case "Addf":
		self.verifyAddfCall(callExpr)
	}

	return self
}

func (self *addfVisitor) verifyAddCall(callExpr *ast.CallExpr) {
	// if .Add is called with more than one argument, this isn't the call site we are looking for
	if len(callExpr.Args) != 1 {
		fmt.Println(`please update the "collector_addf" linter, I found a call to collector.Add with more than one argument`)
		return
	}

	// ensure the argument is a call to fmt.Sprintf
	if !isFmtSprintf(callExpr.Args[0]) {
		return
	}

	// Found a match - report the error
	*self.foundError = true
	workDir := asserts.NoError1(os.Getwd())
	relPath := asserts.NoError1(filepath.Rel(workDir, self.path))
	position := self.fileSet.Position(callExpr.Pos())
	fmt.Printf("%s:%d  Please use the .Addf method to add formatted strings.\n", relPath, position.Line)
}

func (self *addfVisitor) verifyAddfCall(callExpr *ast.CallExpr) {
	// Check if there's at least one argument
	if len(callExpr.Args) > 1 {
		return
	}

	// Found a match - report the error
	*self.foundError = true
	workDir := asserts.NoError1(os.Getwd())
	relPath := asserts.NoError1(filepath.Rel(workDir, self.path))
	position := self.fileSet.Position(callExpr.Pos())
	fmt.Printf("%s:%d: Please use the .Add method since no formatting is happening.\n", relPath, position.Line)
}

func (self *addfVisitor) isCollectorType(expr ast.Expr) bool {
	if self.typeInfo == nil {
		return false
	}
	typ := self.typeInfo.TypeOf(expr)
	if typ == nil {
		return false
	}

	// Get the underlying type name
	typeName := typ.String()

	// Check for both value and pointer receivers
	return strings.Contains(typeName, "stringslice.Collector")
}

func isFmtSprintf(expr ast.Expr) bool {
	callExpr, isCallExpr := expr.(*ast.CallExpr)
	if !isCallExpr {
		return false
	}

	// Check if this is a selector expression (package.Function)
	selectorExpr, isSelectorExpr := callExpr.Fun.(*ast.SelectorExpr)
	if !isSelectorExpr {
		return false
	}

	// Check if the selector is "Sprintf"
	if selectorExpr.Sel.Name != "Sprintf" {
		return false
	}

	// Check if the package identifier is "fmt"
	pkgIdent, isPkgIdent := selectorExpr.X.(*ast.Ident)
	if !isPkgIdent {
		return false
	}

	return pkgIdent.Name == "fmt"
}

func shouldIgnorePath(path string) bool {
	return strings.Contains(path, "vendor/") ||
		strings.Contains(path, ".git/") ||
		strings.HasSuffix(path, "internal/gohacks/stringslice/collector.go")
}
