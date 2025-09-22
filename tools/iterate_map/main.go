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

// append this comment to a line to ignore the problem
const ignoreComment = "// okay to iterate map in random order"

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

			visitor := &mapIterationVisitor{
				fset:     pkg.Fset,
				path:     pkg.GoFiles[i],
				typeInfo: pkg.TypesInfo,
				errors:   &errors,
				file:     file,
			}
			ast.Walk(visitor, file)
		}
	}

	if errors > 0 {
		os.Exit(1)
	}
}

type mapIterationVisitor struct {
	fset     *token.FileSet
	path     string
	typeInfo *types.Info
	errors   *int
	file     *ast.File
}

func (v *mapIterationVisitor) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return v
	}
	if !v.isMapIteration(rangeStmt) {
		return v
	}
	pos := v.fset.Position(rangeStmt.Pos())
	if v.hasIgnoreComment(rangeStmt) {
		return v
	}
	*v.errors++
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return v
	}
	relPath, err := filepath.Rel(workDir, v.path)
	if err != nil {
		relPath = v.path
	}
	fmt.Printf("%s:%d\n", relPath, pos.Line)
	return v
}

func (v *mapIterationVisitor) isMapIteration(rangeStmt *ast.RangeStmt) bool {
	if v.typeInfo == nil {
		return false
	}
	typ := v.typeInfo.TypeOf(rangeStmt.X)
	if typ == nil {
		return false
	}
	return v.isMapType(typ)
}

func (v *mapIterationVisitor) hasIgnoreComment(rangeStmt *ast.RangeStmt) bool {
	if v.file.Comments == nil {
		return false
	}
	rangePos := v.fset.Position(rangeStmt.Pos())
	for _, commentGroup := range v.file.Comments {
		for _, comment := range commentGroup.List {
			commentPos := v.fset.Position(comment.Pos())
			if commentPos.Line == rangePos.Line && strings.HasPrefix(comment.Text, ignoreComment) {
				return true
			}
		}
	}
	return false
}

func (v *mapIterationVisitor) isMapType(typ types.Type) bool {
	switch t := typ.Underlying().(type) {
	case *types.Map:
		return true
	case *types.Pointer:
		return v.isMapType(t.Elem())
	case *types.Named:
		return v.isMapType(t.Underlying())
	}
	return false
}

func shouldIgnorePath(path string) bool {
	return strings.HasPrefix(path, "vendor/") ||
		strings.HasPrefix(path, ".git/") ||
		strings.Contains(path, "testdata/")
}
