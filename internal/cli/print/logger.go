package print

import (
	"fmt"
	"regexp"
	"sync"

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
	logStartOnce.Do(func() {
		logStartRegex = regexp.MustCompile(`%(?:s|d|w)`) // Split by any of "%s", "%d", or "%w"
	})
	parts := logStartRegex.Split(template, -1)
	matches := logStartRegex.FindAllString(template, -1)
	for i := range data {
		fmt.Print(colors.Bold().Styled(parts[i]))
		format := matches[i]
		if format == "%s" {
			fmt.Print(colors.BoldCyan().Styled(fmt.Sprintf(format, data[i])))
		} else {
			fmt.Print(colors.Bold().Styled(fmt.Sprintf(format, data[i])))
		}
	}
	fmt.Print(colors.Bold().Styled(parts[len(parts)-1]))
}

var (
	logStartOnce  sync.Once
	logStartRegex *regexp.Regexp
)

func (self Logger) Success(message string) {
	self.Log(colors.BoldGreen().Styled(message))
}
