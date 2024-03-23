package print

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/cli/format"
)

// LabelAndValue prints the label bolded and underlined
// the value indented on the next line
// followed by an empty line.
func LabelAndValue(label, value string) {
	Header(label)
	fmt.Println(format.Indent(value))
	fmt.Println()
}
