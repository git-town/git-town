package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

// --- Configuration ---
// CHANGE THIS to the actual package import path where your Option[T] is defined.
const targetPackagePath = "github.com/git-town/git-town/v19/pkg/prelude"

const targetTypeName = "Option"

// --- End Configuration ---

var Analyzer = &analysis.Analyzer{
	Name:     "optioncompare",
	Doc:      fmt.Sprintf("Checks for direct == comparisons between Option types"),
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
	Flags:    *flag.NewFlagSet("optioncompare", flag.ExitOnError), // Include flags for configuration
}

func init() {
	// Define the flag for the package path within the analyzer's FlagSet
	// The default value is set above where targetPackagePath is initialized.
	// Analyzer.Flags.StringVar(targetPackagePath, "optionpkg", "main", "Import path of the package defining Option[T]")
}

func main() {
	// The flag is parsed by singlechecker.Main automatically
	singlechecker.Main(Analyzer)
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspectorInstance := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.BinaryExpr)(nil), // We are interested in binary expressions like 'a == b'
	}

	inspectorInstance.Preorder(nodeFilter, func(n ast.Node) {
		binExpr := n.(*ast.BinaryExpr)

		// 1. Check if the operator is ==
		if binExpr.Op != token.EQL {
			return
		}

		// 2. Get type information for left and right operands
		leftTypeInfo := pass.TypesInfo.TypeOf(binExpr.X)
		rightTypeInfo := pass.TypesInfo.TypeOf(binExpr.Y)

		if leftTypeInfo == nil || rightTypeInfo == nil {
			// Skip if type info is unavailable for some reason
			return
		}

		// 3. Allow comparisons with untyped or typed nil
		// If one side is nil, it's okay (e.g., opt == nil)
		// if isNil(pass, binExpr.X) || isNil(pass, binExpr.Y) {
		// 	return
		// }

		// 4. Check if *both* operands are the target Option[T] type
		if isTargetOptionType(leftTypeInfo) && isTargetOptionType(rightTypeInfo) {
			// Construct a user-friendly type name (e.g., main.Option[int])
			typeName := leftTypeInfo.String()
			// Sometimes type String() might include the package path redundantly, clean it up if needed for display
			typeName = strings.Replace(typeName, targetPackagePath+".", "", 1) // Basic cleanup

			pass.Reportf(binExpr.OpPos, // Report at the position of the '==' operator
				"direct comparison of type %s using ==; use the Equal() method instead", typeName)
		}
	})

	return nil, nil
}

// isTargetOptionType checks if a type is an instance of the target Option[T] generic type
// from the configured package.
func isTargetOptionType(typ types.Type) bool {
	if typ == nil {
		return false
	}

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
	if named.Obj() == nil || named.Obj().Name() != targetTypeName {
		return false
	}

	// Check the package path where the type is defined
	pkg := named.Obj().Pkg()
	if pkg == nil || pkg.Path() != targetPackagePath {
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
	if origin.Obj() == nil || origin.Obj().Name() != targetTypeName ||
		origin.Obj().Pkg() == nil || origin.Obj().Pkg().Path() != targetPackagePath {
		return false
	}

	// Confirm the origin type actually has type parameters (is generic)
	if origin.TypeParams().Len() == 0 {
		return false // The original type wasn't generic
	}

	// If all checks pass, it's an instantiation of our target Option[T] type.
	return true
}

// isNil checks if an expression resolves to the untyped nil literal or a typed nil value.
func isNil(pass *analysis.Pass, expr ast.Expr) bool {
	// Check for untyped nil literal
	if basic, ok := pass.TypesInfo.TypeOf(expr).(*types.Basic); ok && basic.Kind() == types.UntypedNil {
		return true
	}
	// Check for typed nil identifier
	if id, ok := expr.(*ast.Ident); ok && id.Name == "nil" {
		// Check if it's the predefined nil object
		obj := pass.TypesInfo.ObjectOf(id)
		return obj != nil && obj.Name() == "nil" && obj.Parent() == types.Universe // Ensure it's the universal nil
	}
	return false
}
