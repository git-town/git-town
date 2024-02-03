package print

import (
	"fmt"

	"github.com/muesli/termenv"
)

// Error prints the given error message to the console.
func Error(err error) {
	boldRed := termenv.String().Bold().Foreground(termenv.ANSIRed)
	fmt.Println(boldRed.Styled("\nError: " + err.Error() + "\n"))
}
