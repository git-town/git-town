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

// file paths to ignore
var ignorePaths = []string{ //nolint:gochecknoglobals
	"internal/gohacks/slice/hoist_test.go",
	"internal/gohacks/slice/remove_at_test.go",
	"internal/gohacks/slice/remove_test.go",
	"internal/gohacks/slice/truncate_last_test.go",
	"internal/gohacks/stringslice/collector_test.go",
	"internal/gohacks/stringslice/lines_with_prefix_test.go",
	"internal/hosting/detect_test.go",
	"internal/undo/undobranches/branch_changes_test.go",
	"internal/undo/undobranches/branch_span_test.go",
	"internal/undo/undoconfig/config_test.go",
	"internal/vm/optimizer/remove_duplicate_checkout_test.go",
	"internal/vm/program/program_test.go",
	"pkg/prelude/option_test.go",
	"test/datatable/branch_setup_table_test.go",
	"test/datatable/data_table_test.go",
	"test/filesystem/copy_directory_test.go",
	"test/fixture/factory_test.go",
	"test/fixture/fixture_test.go",
	"test/output/git_town_output_parser_test.go",
	"test/subshell/test_runner_test.go",
	"test/testgit/locations_test.go",
	"test/testruntime/testruntime_test.go",
	"tools/format_unittests/format_unittests_test.go",
	// Not owned.
	"vendor/",
}

func main() {
	issues := Issues{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".go") || isIgnoredPath(path) {
			return err
		}
		fileIssues, err := LintFile(path, nil)
		if err != nil {
			return err
		}
		issues = append(issues, fileIssues...)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(issues) > 0 {
		fmt.Println(issues)
	}
	os.Exit(len(issues))
}

// LintFile lints file specified by a path.
// If contents is not nil the path is only used for positional information.
// If contents is nil the file is loaded from path.
func LintFile(path string, contents any) (Issues, error) {
	result := Issues{}
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, contents, 0)
	if err != nil {
		return result, err
	}
	ast.Inspect(file, func(node ast.Node) bool {
		var nodeIssues Issues
		var nodeErr error
		// To lint anonymous functions add `case *ast.FuncLit`.
		// This would enforce subtest sorting recursively.
		switch typedNode := node.(type) { //nolint:gocritic
		// Selects top-level function declarations.
		case *ast.FuncDecl:
			nodeIssues, nodeErr = lintFuncDecl(typedNode, fileSet)
		}
		if nodeErr != nil {
			// Propagate nodeErr to the parent scope.
			err = nodeErr
			return false
		}
		result = append(result, nodeIssues...)
		return true
	})
	return result, err
}

func isIgnoredPath(path string) bool {
	for _, ignore := range ignorePaths {
		if strings.HasPrefix(path, ignore) {
			return true
		}
	}
	return false
}
