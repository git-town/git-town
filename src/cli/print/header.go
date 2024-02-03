package print

import (
	"fmt"

	"github.com/muesli/termenv"
)

func Header(text string) {
	boldUnderline := termenv.String().Bold().Underline()
	fmt.Println(boldUnderline.Styled(text + ":"))
}
