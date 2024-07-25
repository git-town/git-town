package main_test

import (
	"os"
	"testing"

	structsSorted "github.com/git-town/git-town/tools/structs_sorted"
	"github.com/shoenig/test/must"
)

const testPath = "test.go"

func TestStructsSorted(t *testing.T) {
	t.Parallel()

	t.Run("LintFile", func(t *testing.T) {
		t.Parallel()
		t.Run("unsorted definition", func(t *testing.T) {
			give := `
package main
type MyStruct struct {
	field2 int // this field should not be first
	field1 int // this field should be first
}
`
			createTestFile(give)
			defer os.Remove(testPath)
			have := structsSorted.LintFile(testPath).String()
			want := `
test.go:3:6 unsorted fields, expected order:

field1
field2

`[1:]
			must.EqOp(t, want, have)
		})

		t.Run("definition without fields", func(t *testing.T) {
			give := `
package main
type MyStruct struct {}
`
			createTestFile(give)
			defer os.Remove(testPath)
			have := structsSorted.LintFile(testPath).String()
			want := ""
			must.EqOp(t, want, have)
		})

		t.Run("ignored definition", func(t *testing.T) {
			give := `
package main
type Change struct {
	field2 int
	field1 int
}
`
			createTestFile(give)
			defer os.Remove(testPath)
			have := structsSorted.LintFile(testPath).String()
			want := ""
			must.EqOp(t, want, have)
		})

		t.Run("unsorted instantiation", func(t *testing.T) {
			give := `
package main
type MyStruct struct {
	field1 int
	field2 int
}
func main() {
	foo := MyStruct{
		field2: 2,
		field1: 1,
	}
}
`
			createTestFile(give)
			defer os.Remove(testPath)
			have := structsSorted.LintFile(testPath).String()
			want := `
test.go:8:9 unsorted fields, expected order:

field1
field2

`[1:]
			must.EqOp(t, want, have)
		})

		t.Run("instantiation without fields", func(t *testing.T) {
			give := `
package main
type MyStruct struct {}
func main() {
	foo := MyStruct{}
}
`
			createTestFile(give)
			defer os.Remove(testPath)
			have := structsSorted.LintFile(testPath).String()
			want := ""
			must.EqOp(t, want, have)
		})

		t.Run("ignored instantiation", func(t *testing.T) {
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
			createTestFile(give)
			defer os.Remove(testPath)
			have := structsSorted.LintFile(testPath).String()
			want := ""
			must.EqOp(t, want, have)
		})
	})
}

func createTestFile(text string) {
	file := os.WriteFile(testPath, []byte(text), 0o600)
	if file != nil {
		panic(file.Error())
	}
}
