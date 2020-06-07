package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetUserInput reads input from the user and returns it.
func GetUserInput() string {
	inputReader := bufio.NewReader(os.Stdin)
	text, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error getting user input: %v", err)
		os.Exit(1)
	}
	return strings.TrimSpace(text)
}
