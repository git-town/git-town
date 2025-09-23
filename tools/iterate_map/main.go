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

func (self *mapIterationVisitor) Visit(node ast.Node) ast.Visitor {
	rangeStmt, isRangeStmt := node.(*ast.RangeStmt)
	if !isRangeStmt {
		return self
	}
	if !self.isMapIteration(rangeStmt) {
		return self
	}
	if self.hasIgnoreComment(rangeStmt) {
		return self
	}
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
	position := self.fset.Position(rangeStmt.Pos())
	fmt.Printf("%s:%d\n", relPath, position.Line)
	return self
}

func (self *mapIterationVisitor) hasIgnoreComment(rangeStmt *ast.RangeStmt) bool {
	if self.file.Comments == nil {
		return false
	}
	rangePos := self.fset.Position(rangeStmt.Pos())
	for _, commentGroup := range self.file.Comments {
		for _, comment := range commentGroup.List {
			commentPos := self.fset.Position(comment.Pos())
			if commentPos.Line == rangePos.Line && strings.HasPrefix(comment.Text, ignoreComment) {
				return true
			}
		}
	}
	return false
}

func (self *mapIterationVisitor) isMapIteration(rangeStmt *ast.RangeStmt) bool {
	if self.typeInfo == nil {
		return false
	}
	typ := self.typeInfo.TypeOf(rangeStmt.X)
	if typ == nil {
		return false
	}
	return self.isMapType(typ)
}

func (self *mapIterationVisitor) isMapType(typ types.Type) bool {
	switch t := typ.Underlying().(type) {
	case *types.Map:
		return true
	case *types.Pointer:
		return self.isMapType(t.Elem())
	case *types.Named:
		return self.isMapType(t.Underlying())
	}
	return false
}

func shouldIgnorePath(path string) bool {
	return strings.HasPrefix(path, "vendor/") ||
		strings.HasPrefix(path, ".git/")
}
