package drivers

import (
	"fmt"

	"github.com/fatih/color"
)

func printLog(message string) {
	fmt.Println()
	color.New(color.Bold).Println(message)
}
