package opcodes

import (
	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// BrowserOpen displays the existing proposal with the given URL in the browser.
type BrowserOpen struct {
	URL string
}

func (self *BrowserOpen) Run(args shared.RunArgs) error {
	browser.Open(self.URL, args.Frontend)
	return nil
}
