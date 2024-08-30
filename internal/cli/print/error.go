package print

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/colors"
)

// Error prints the given error message to the console.
func Error(err error) {
	fmt.Println(colors.BoldRed().Styled("\nError: " + err.Error() + "\n"))
}
