package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

var (
	// file paths to ignore
	ignorePaths = []string{ //nolint:gochecknoglobals
		"vendor/",
		"tools/structs_sorted/test.go",
	}

	// struct types to ignore
	ignoreTypes = []string{ //nolint:gochecknoglobals
		"BranchSpan",
		"Change",
		"InconsistentChange",
		"Parts",
		"ProdRunner",
	}
)

type issue struct {
	expected   []string
	pos        token.Position
	structName string
}

func (self issue) String() string {
	return fmt.Sprintf(
		"%s:%d:%d unsorted fields in %s. Expected order:\n\n%s\n\n",
		self.pos.Filename,
		self.pos.Line,
		self.pos.Column,
		self.structName,
		strings.Join(self.expected, "\n"),
	)
}

type issues []issue

func (self issues) String() string {
	result := strings.Builder{}
	for _, issue := range self {
		result.WriteString(issue.String())
	}
	return result.String()
}

func main() {
	switch {
	case len(os.Args) == 1 || len(os.Args) > 2:
		displayUsage()
	case len(os.Args) == 2 && os.Args[1] == "run":
		lintFiles()
	case len(os.Args) == 2 && os.Args[1] == "test":
		runTests()
	default:
		fmt.Printf("Error: unknown argument: %s", os.Args[1])
		os.Exit(1)
	}
}

func displayUsage() {
	fmt.Println(`
This tool lints Go code for alphabetic sorting of struct fields.

Usage: lint <command>

Available commands:
   run   Lints all files
   test  Verifies that this tool works
`[1:])
}

func lintFiles() {
	issues := []issue{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".go") || isIgnored(path) {
			return err
		}
		issues = append(issues, lintFile(path)...)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	printIssues(issues)
	os.Exit(len(issues))
}

func printIssues(issues issues) {
	for _, issue := range issues {
		fmt.Println(issue.String())
	}
}

func lintFile(path string) issues {
	result := issues{}
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return result
	}
	ast.Inspect(file, func(node ast.Node) bool {
		result = append(result, lintStructDefinitions(node, fileSet)...)
		result = append(result, lintStructLiteral(node, fileSet)...)
		return true
	})
	return result
}

func lintStructDefinitions(node ast.Node, fileSet *token.FileSet) issues {
	typeSpec, ok := node.(*ast.TypeSpec)
	if !ok {
		return issues{}
	}
	structName := typeSpec.Name.Name
	if slices.Contains(ignoreTypes, structName) {
		return issues{}
	}
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return issues{}
	}
	fields := structDefFieldNames(structType)
	sortedFields := make([]string, len(fields))
	copy(sortedFields, fields)
	slices.Sort(sortedFields)
	if reflect.DeepEqual(fields, sortedFields) {
		return issues{}
	}
	return issues{
		issue{
			pos:        fileSet.Position(node.Pos()),
			structName: structName,
			expected:   sortedFields,
		},
	}
}

func lintStructLiteral(node ast.Node, fileSet *token.FileSet) issues {
	compositeLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return issues{}
	}
	structType, ok := compositeLit.Type.(*ast.Ident)
	if !ok {
		return issues{}
	}
	var fieldNames []string
	for _, elt := range compositeLit.Elts {
		kvExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		fieldName := kvExpr.Key.(*ast.Ident).Name
		fieldNames = append(fieldNames, fieldName)
	}
	sortedFields := make([]string, len(fieldNames))
	copy(sortedFields, fieldNames)
	slices.Sort(sortedFields)
	if reflect.DeepEqual(fieldNames, sortedFields) {
		return issues{}
	}
	return issues{
		issue{
			expected:   sortedFields,
			pos:        fileSet.Position(node.Pos()),
			structName: structType.Name,
		},
	}
}

func isIgnored(path string) bool {
	for _, ignore := range ignorePaths {
		if strings.HasPrefix(path, ignore) {
			return true
		}
	}
	return false
}

func structDefFieldNames(structType *ast.StructType) []string {
	var result []string
	for _, field := range structType.Fields.List {
		if field.Names != nil {
			result = append(result, field.Names[0].Name)
		}
	}
	return result
}

func structInstFieldNames(compositeLit *ast.CompositeLit) []string {
	var result []string
	for _, expr := range compositeLit.Elts {
		if kvExpr, ok := expr.(*ast.KeyValueExpr); ok {
			if ident, ok := kvExpr.Key.(*ast.Ident); ok {
				result = append(result, ident.Name)
			}
		}
	}
	return result
}

/************************************************************************************
 * TESTS
 */

func runTests() {
	testUnsortedDefinition()
	testUnsortedCall()
	fmt.Println()
}

func testUnsortedDefinition() {
	give := `
package main

type Unsorted struct {
	field2 int // this field should not be first
	field1 int // this field should be first
}`
	path := "test.go"
	file := os.WriteFile(path, []byte(give), 0644)
	if file != nil {
		panic(file.Error())
	}
	defer os.Remove(path)
	have := lintFile(path).String()
	want := `
test.go:4:6 unsorted fields in Unsorted. Expected order:

field1
field2

`[1:]
	assertEqual(want, have, "testUnsortedDefinition")
}

func testUnsortedCall() {
	give := `
package main

type Foo struct {
	field1 int
	field2 int
}

func main() {
	foo := Foo{
		field2: 2,
		field1: 1,
	}
}
`
	path := "test.go"
	file := os.WriteFile(path, []byte(give), 0644)
	if file != nil {
		panic(file.Error())
	}
	defer os.Remove(path)
	have := lintFile(path).String()
	want := `
test.go:10:9 unsorted fields in Foo. Expected order:

field1
field2

`[1:]
	assertEqual(want, have, "testUnsortedDefinition")
}

func assertEqual[T comparable](want, have T, testName string) {
	fmt.Print(".")
	if have != want {
		fmt.Printf("\nTEST FAILURE in %q\n", testName)
		fmt.Println("\n\nWANT")
		fmt.Println("--------------------------------------------------------")
		fmt.Println(want)
		fmt.Println("\n\nHAVE")
		fmt.Println("--------------------------------------------------------")
		fmt.Println(have)
		os.Exit(1)
	}
}
