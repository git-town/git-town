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
	expected   []string       // the expected order of fields
	pos        token.Position // file, line, and column of the issue
	structName string         // name of the struct that has the problem described by this issue
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
   run   Lints all files Go in the current directory and subdirectories
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
	fmt.Println(issues)
	os.Exit(len(issues))
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
	if len(fields) == 0 {
		return issues{}
	}
	sortedFields := make([]string, len(fields))
	copy(sortedFields, fields)
	slices.Sort(sortedFields)
	if reflect.DeepEqual(fields, sortedFields) {
		return issues{}
	}
	return issues{
		issue{
			expected:   sortedFields,
			pos:        fileSet.Position(node.Pos()),
			structName: structName,
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
	structName := structType.Name
	if slices.Contains(ignoreTypes, structName) {
		return issues{}
	}
	fieldNames := make([]string, 0, len(compositeLit.Elts))
	for _, elt := range compositeLit.Elts {
		kvExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kvExpr.Key.(*ast.Ident)
		if !ok {
			continue
		}
		fieldName := ident.Name
		fieldNames = append(fieldNames, fieldName)
	}
	if len(fieldNames) == 0 {
		return issues{}
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
			structName: structName,
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

/************************************************************************************
 * TESTS
 */

func runTests() {
	testUnsortedDefinition()
	testDefinitionWithoutFields()
	testIgnoredDefinition()
	testUnsortedInstantiation()
	testInstantiationWithoutFields()
	testIgnoredInstantiation()
	fmt.Println()
}

func testUnsortedDefinition() {
	give := `
package main
type Unsorted struct {
	field2 int // this field should not be first
	field1 int // this field should be first
}
`
	path := "test.go"
	file := os.WriteFile(path, []byte(give), 0o600)
	if file != nil {
		panic(file.Error())
	}
	defer os.Remove(path)
	have := lintFile(path).String()
	want := `
test.go:3:6 unsorted fields in Unsorted. Expected order:

field1
field2

`[1:]
	assertEqual(want, have, "testUnsortedDefinition")
}

func testDefinitionWithoutFields() {
	give := `
package main
type Foo struct {}`
	path := "test.go"
	file := os.WriteFile(path, []byte(give), 0o600)
	if file != nil {
		panic(file.Error())
	}
	defer os.Remove(path)
	have := lintFile(path).String()
	want := ""
	assertEqual(want, have, "testDefinitionWithoutFields")
}

func testIgnoredDefinition() {
	give := `
package main
type Change struct {
	field2 int
	field1 int
}
`
	path := "test.go"
	file := os.WriteFile(path, []byte(give), 0o600)
	if file != nil {
		panic(file.Error())
	}
	defer os.Remove(path)
	have := lintFile(path).String()
	want := ""
	assertEqual(want, have, "testIgnoredDefinition")
}

func testUnsortedInstantiation() {
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
	file := os.WriteFile(path, []byte(give), 0o600)
	if file != nil {
		panic(file.Error())
	}
	defer os.Remove(path)
	have := lintFile(path).String()
	want := `
test.go:8:9 unsorted fields in Foo. Expected order:

field1
field2

`[1:]
	assertEqual(want, have, "testUnsortedInstantiation")
}

func testInstantiationWithoutFields() {
	give := `
package main
type Foo struct {
}
func main() {
	foo := Foo{}
}
`
	path := "test.go"
	file := os.WriteFile(path, []byte(give), 0o600)
	if file != nil {
		panic(file.Error())
	}
	defer os.Remove(path)
	have := lintFile(path).String()
	want := ""
	assertEqual(want, have, "testInstantiationWithoutFields")
}

func testIgnoredInstantiation() {
	give := `
package main
type Change struct {
	field1 int
	field2 int
}
func main() {
	foo := Change{
		field2: 2,
		field1: 1,
	}
}
`
	path := "test.go"
	file := os.WriteFile(path, []byte(give), 0o600)
	if file != nil {
		panic(file.Error())
	}
	defer os.Remove(path)
	have := lintFile(path).String()
	want := ""
	assertEqual(want, have, "testIgnoredInstantiation")
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
