package drivers

import (
	"fmt"

	"github.com/Originate/exit"
	"github.com/fatih/color"
)

func printLog(message string) {
	fmt.Println()
	_, err := color.New(color.Bold).Println(message)
	exit.If(err)
}
