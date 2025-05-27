package print

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/colors"
)

func Header(text string) {
	boldUnderline := colors.BoldUnderline()
	fmt.Println(boldUnderline.Styled(text + ":"))
}
