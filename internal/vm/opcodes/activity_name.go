package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v20/internal/cli/colors"
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// ActivityName displays the name of the current high-level activity
type ActivityName struct {
	Text                    string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ActivityName) Run(args shared.RunArgs) error {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println(colors.Bold().Styled(self.Text))
	return nil
}
