package opcodes

import (
	"github.com/git-town/git-town/v19/internal/browser"
	"github.com/git-town/git-town/v19/internal/vm/shared"
)

// ProposalShow displays the existing proposal with the given URL in the browser.
type ProposalShow struct {
	URL                     string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalShow) Run(args shared.RunArgs) error {
	browser.Open(self.URL, args.Frontend, args.Backend)
	return nil
}
