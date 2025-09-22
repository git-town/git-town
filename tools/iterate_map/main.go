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

// append a comment starting with this text
// to a line that this linter warns about
// to silence the warning
const ignoreComment = "// okay to iterate the map in random order"

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
				errors:   &errors,
				file:     file,
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

type mapIterationVisitor struct {
	errors   *int
	file     *ast.File
	fset     *token.FileSet
	path     string
	typeInfo *types.Info
}

func (visitor *mapIterationVisitor) Visit(node ast.Node) ast.Visitor {
	rangeStmt, isRangeStmt := node.(*ast.RangeStmt)
	if !isRangeStmt {
		return visitor
	}
	if !visitor.isMapIteration(rangeStmt) {
		return visitor
	}
	if visitor.hasIgnoreComment(rangeStmt) {
		return visitor
	}
	*visitor.errors++
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return visitor
	}
	relPath, err := filepath.Rel(workDir, visitor.path)
	if err != nil {
		fmt.Println(err.Error())
		return visitor
	}
	position := visitor.fset.Position(rangeStmt.Pos())
	fmt.Printf("%s:%d\n", relPath, position.Line)
	return visitor
}

func (visitor *mapIterationVisitor) isMapIteration(rangeStmt *ast.RangeStmt) bool {
	if visitor.typeInfo == nil {
		return false
	}
	typ := visitor.typeInfo.TypeOf(rangeStmt.X)
	if typ == nil {
		return false
	}
	return visitor.isMapType(typ)
}

func (visitor *mapIterationVisitor) hasIgnoreComment(rangeStmt *ast.RangeStmt) bool {
	if visitor.file.Comments == nil {
		return false
	}
	rangePos := visitor.fset.Position(rangeStmt.Pos())
	for _, commentGroup := range visitor.file.Comments {
		for _, comment := range commentGroup.List {
			commentPos := visitor.fset.Position(comment.Pos())
			if commentPos.Line == rangePos.Line && strings.HasPrefix(comment.Text, ignoreComment) {
				return true
			}
		}
	}
	return false
}

func (visitor *mapIterationVisitor) isMapType(typ types.Type) bool {
	switch t := typ.Underlying().(type) {
	case *types.Map:
		return true
	case *types.Pointer:
		return visitor.isMapType(t.Elem())
	case *types.Named:
		return visitor.isMapType(t.Underlying())
	}
	return false
}

func shouldIgnorePath(path string) bool {
	return strings.HasPrefix(path, "vendor/") ||
		strings.HasPrefix(path, ".git/")
}
