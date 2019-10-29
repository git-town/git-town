package drivers

import (
	"fmt"

	"github.com/fatih/color"
)

func printLog(message string) {
	fmt.Println()
	_, err := color.New(color.Bold).Println(message)
	if err != nil {
		panic(err)
	}
}
