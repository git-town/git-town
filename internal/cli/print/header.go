package print

import (
	"fmt"

	"github.com/git-town/git-town/v21/pkg/colors"
)

func Header(text string) {
	boldUnderline := colors.BoldUnderline()
	fmt.Println(boldUnderline.Styled(text + ":"))
}
