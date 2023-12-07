package print

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/format"
	"github.com/git-town/git-town/v11/src/cli/io"
)

// LabelAndValue prints the label bolded and underlined
// the value indented on the next line
// followed by an empty line.
func LabelAndValue(label, value string) {
	Header(label)
	io.Println(format.Indent(value))
	fmt.Println()
}
