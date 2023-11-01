package print

import "github.com/git-town/git-town/v10/src/cli/io"

func Entry(label, value string) {
	print()
	io.Printf("  %s: %s\n", label, value)
}
