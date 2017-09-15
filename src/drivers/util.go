package drivers

import (
	"fmt"

	"github.com/Originate/git-town/src/exit"
	"github.com/fatih/color"
)

func printLog(message string) {
	fmt.Println()
	_, err := color.New(color.Bold).Println(message)
	exit.On(err)
}
