package print

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/pkg/colors"
)

// The Logger logger logs activities of a particular component on the CLI.
type Logger struct{}

func (self Logger) Failed(failure string) {
	self.Log(colors.BoldRed().Styled(fmt.Sprintf("%v\n", failure)))
}

func (self Logger) Finished(err error) {
	if err != nil {
		self.Failed(err.Error())
	} else {
		self.Ok()
	}
}

func (self Logger) Log(text string) {
	fmt.Println(text)
}

func (self Logger) Ok() {
	self.Success("ok")
}

func (self Logger) Start(template string, data ...any) {
	fmt.Println()
	if len(data) == 0 {
		fmt.Print(colors.Bold().Styled(template))
		return
	}
	parts := strings.Split(template, "%s")
	for i, part := range parts {
		fmt.Print(colors.Bold().Styled(part))
		fmt.Print(colors.BoldCyan().Styled(fmt.Sprintf("%s", data[i])))
	}
	fmt.Print(colors.Bold().Styled(parts[len(parts)-1]))
}

func (self Logger) Success(message string) {
	self.Log(colors.BoldGreen().Styled(message))
}
