package print

import "github.com/git-town/git-town/v9/src/cli/io"

func Entry(label, value string) {
	print()
	io.Printf("  %s: %s\n", label, value)
}
