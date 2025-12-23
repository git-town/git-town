package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{
		Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		fmt.Printf("ERROR loading packages: %s\n", err)
		os.Exit(1)
	}
	errors := 0
	for _, pkg := range pkgs {
		for i, file := range pkg.Syntax {
			if shouldIgnorePath(pkg.GoFiles[i]) {
				continue
			}
			visitor := &fmtSprintfVisitor{
				errors:   &errors,
				fset:     pkg.Fset,
				path:     pkg.GoFiles[i],
				typeInfo: pkg.TypesInfo,
			}
			ast.Walk(visitor, file)
		}
	}
	if errors > 0 {
		os.Exit(1)
	}
}

type fmtSprintfVisitor struct {
	errors   *int
	fset     *token.FileSet
	path     string
	typeInfo *types.Info
}

func (self *fmtSprintfVisitor) Visit(node ast.Node) ast.Visitor {
	callExpr, isCallExpr := node.(*ast.CallExpr)
	if !isCallExpr {
		return self
	}

	// Check if this is a method call (selector expression)
	selectorExpr, isSelectorExpr := callExpr.Fun.(*ast.SelectorExpr)
	if !isSelectorExpr {
		return self
	}

	// Check if the method name is "Add"
	if selectorExpr.Sel.Name != "Add" {
		return self
	}

	// Check if the receiver is of type stringslice.Collector
	if !self.isCollectorType(selectorExpr.X) {
		return self
	}

	// Check if there's exactly one argument and it's a call to fmt.Sprintf
	if len(callExpr.Args) != 1 {
		return self
	}

	if !self.isFmtSprintf(callExpr.Args[0]) {
		return self
	}

	// Found a match - report the error
	*self.errors++
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return self
	}
	relPath, err := filepath.Rel(workDir, self.path)
	if err != nil {
		fmt.Println(err.Error())
		return self
	}
	position := self.fset.Position(callExpr.Pos())
	fmt.Printf("%s:%d: Please use the .AddF method to add formatted strings.\n", relPath, position.Line)
	return self
}

func (self *fmtSprintfVisitor) isCollectorType(expr ast.Expr) bool {
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

func (self *fmtSprintfVisitor) isFmtSprintf(expr ast.Expr) bool {
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
	return strings.HasPrefix(path, "vendor/") ||
		strings.HasPrefix(path, ".git/") ||
		strings.HasSuffix(path, "internal/gohacks/stringslice/collector.go")
}
