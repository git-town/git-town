package drivers

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func printLog(message string) {
	fmt.Println()
	_, err := color.New(color.Bold).Println(message)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
