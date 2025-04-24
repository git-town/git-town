package opcodes

import (
	"github.com/git-town/git-town/v19/internal/browser"
	"github.com/git-town/git-town/v19/internal/vm/shared"
)

// BrowserOpen opens the browser with the given URL.
type BrowserOpen struct {
	URL                     string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BrowserOpen) Run(args shared.RunArgs) error {
	browser.Open(self.URL, args.Frontend, args.Backend)
	return nil
}
