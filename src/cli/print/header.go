package print

import (
	"github.com/fatih/color"
	"github.com/git-town/git-town/v11/src/cli/io"
)

func Header(text string) {
	boldUnderline := color.New(color.Bold).Add(color.Underline)
	io.PrintlnColor(boldUnderline, text+":")
}
