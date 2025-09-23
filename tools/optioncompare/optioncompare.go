package main

import (
	"flag"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	packagePath = "github.com/git-town/git-town/v22/pkg/prelude"
	typeName    = "Option"
)

func main() {
	singlechecker.Main(&analysis.Analyzer{
		Name:     "optioncompare",
		Doc:      "Ensures no == comparison between Option types",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      execute,
		Flags:    *flag.NewFlagSet("optioncompare", flag.ExitOnError),
	})
}

func execute(pass *analysis.Pass) (any, error) { //nolint:ireturn
	inspectorInstance := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.BinaryExpr)(nil), // We are interested in binary expressions like 'a == b'
	}
	inspectorInstance.Preorder(nodeFilter, func(node ast.Node) {
		binExpr := node.(*ast.BinaryExpr)
		// operator must be ==
		if binExpr.Op != token.EQL {
			return
		}
		// get types for left and right operands
		leftTypeInfo := pass.TypesInfo.TypeOf(binExpr.X)
		rightTypeInfo := pass.TypesInfo.TypeOf(binExpr.Y)
		if leftTypeInfo == nil || rightTypeInfo == nil {
			panic("cannot determine types")
		}
		// both operands must be an Option type
		if !isOptionType(leftTypeInfo) || !isOptionType(rightTypeInfo) {
			return
		}
		// bingo --> report the problem
		pass.Reportf(binExpr.OpPos, "must compare Options using .Equal instead of ==")
	})
	return nil, nil
}

// isOptionType indicates whether the given type is an Option[T] generic type
func isOptionType(typ types.Type) bool {
	// Check if it's a Named type (like main.Option[int])
	// Using Underlying() can help resolve type aliases if needed, but start with direct type.
	named, ok := typ.(*types.Named)
	if !ok {
		// It might be a pointer *Option[T]. Let's check the element type.
		// Direct comparison of pointers (&opt1 == &opt2) is identity comparison, often valid.
		// Comparing the value *opt1 == opt2 requires *opt1 to be checked.
		// Let's focus on the direct opt1 == opt2 case as requested initially.
		// If you need to lint *opt1 == *opt2, you'd need pointer handling here.
		return false
	}

	// Check the type name itself
	if named.Obj() == nil || named.Obj().Name() != typeName {
		return false
	}

	// Check the package path where the type is defined
	pkg := named.Obj().Pkg()
	if pkg == nil || pkg.Path() != packagePath {
		// This check prevents flagging types named "Option" from other packages.
		// It also handles vendored paths correctly if the paths match.
		return false
	}

	// Verify it's an instantiation of the generic type by checking its origin.
	// The Origin() method returns the generic type definition (`Option[T]`)
	// from an instantiated type (`Option[int]`).
	origin := named.Origin()
	if origin == nil {
		// If Origin() is nil, it might be the generic type definition itself,
		// or a non-generic named type. We are interested in instantiations.
		return false
	}

	// Double-check the origin's name and package for robustness
	if origin.Obj() == nil ||
		origin.Obj().Name() != typeName ||
		origin.Obj().Pkg() == nil ||
		origin.Obj().Pkg().Path() != packagePath {
		return false
	}

	// Confirm the origin type actually has type parameters (is generic)
	if origin.TypeParams().Len() == 0 {
		return false // The original type wasn't generic
	}

	// If all checks pass, it's an instantiation of our target Option[T] type.
	return true
}
