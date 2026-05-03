package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/browser"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BrowserOpen displays the existing proposal with the given URL in the browser.
type BrowserOpen struct {
	URL string
}

func (self *BrowserOpen) Run(args shared.RunArgs) error {
	if args.Config.Value.NormalConfig.BrowserEnabled {
		browser.Open(self.URL, args.Frontend, args.Config.Value.NormalConfig.BrowserExecutable)
	} else {
		fmt.Printf(messages.BrowserOpen, self.URL)
	}
	return nil
}
