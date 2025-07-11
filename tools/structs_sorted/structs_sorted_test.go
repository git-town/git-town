package main_test

import (
	"os"
	"testing"

	structsSorted "github.com/git-town/git-town/tools/structs_sorted"
	"github.com/shoenig/test/must"
)

func TestStructsSorted(t *testing.T) {
	t.Parallel()

	t.Run("LintFile", func(t *testing.T) {
		t.Parallel()
		t.Run("unsorted definition", func(t *testing.T) {
			t.Parallel()
			give := `
package main
type MyStruct struct {
	field2 int // this field should not be first
	field1 int // this field should be first
}
`
			createTestFile(give, "test1.go")
			defer os.Remove("test1.go")
			have := structsSorted.LintFile("test1.go").String()
			want := `
test1.go:3:6 unsorted fields, expected order:

field1
field2

`[1:]
			must.EqOp(t, want, have)
		})

		t.Run("definition without fields", func(t *testing.T) {
			t.Parallel()
			give := `
package main
type MyStruct struct {}
`
			createTestFile(give, "test2.go")
			defer os.Remove("test2.go")
			have := structsSorted.LintFile("test2.go").String()
			want := ""
			must.EqOp(t, want, have)
		})

		t.Run("ignored definition", func(t *testing.T) {
			t.Parallel()
			give := `
package main
type Change struct {
	field2 int
	field1 int
}
`
			createTestFile(give, "test3.go")
			defer os.Remove("test3.go")
			have := structsSorted.LintFile("test3.go").String()
			want := ""
			must.EqOp(t, want, have)
		})

		t.Run("unsorted instantiation", func(t *testing.T) {
			t.Parallel()
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
			createTestFile(give, "test4.go")
			defer os.Remove("test4.go")
			have := structsSorted.LintFile("test4.go").String()
			want := `
test4.go:8:9 unsorted fields, expected order:

field1
field2

`[1:]
			must.EqOp(t, want, have)
		})

		t.Run("instantiation without fields", func(t *testing.T) {
			t.Parallel()
			give := `
package main
type MyStruct struct {}
func main() {
	foo := MyStruct{}
}
`
			createTestFile(give, "test5.go")
			defer os.Remove("test5.go")
			have := structsSorted.LintFile("test5.go").String()
			want := ""
			must.EqOp(t, want, have)
		})

		t.Run("ignored instantiation", func(t *testing.T) {
			t.Parallel()
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
			createTestFile(give, "test6.go")
			defer os.Remove("test6.go")
			have := structsSorted.LintFile("test6.go").String()
			want := ""
			must.EqOp(t, want, have)
		})
	})
}

func createTestFile(text, filename string) {
	file := os.WriteFile(filename, []byte(text), 0o600)
	if file != nil {
		panic(file.Error())
	}
}
