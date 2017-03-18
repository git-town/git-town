package util

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func DoesCommandOuputContain(cmd []string, value string) bool {
	return strings.Contains(GetCommandOutput(cmd...), value)
}

func DoesCommandOuputContainLine(cmd []string, value string) bool {
	list := strings.Split(GetCommandOutput(cmd...), "\n")
	return DoesStringArrayContain(list, value)
}

func DoesStringArrayContain(list []string, value string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

func ExitWithErrorMessage(message string) {
	PrintError(message)
	os.Exit(1)
}

func GetCommandOutput(cmd ...string) string {
	subProcess := exec.Command(cmd[0], cmd[1:]...)
	output, err := subProcess.CombinedOutput()
	if err != nil {
		log.Fatal("Command: ", strings.Join(cmd, " "), "\nOutput: "+string(output), "\nError: ", err)
	}
	return strings.TrimSpace(string(output))
}

func GetUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		log.Fatal("Error getting user input:", err)
	}
	return strings.TrimSpace(text)
}

func PrintError(message string) {
	errHeaderFmt := color.New(color.Bold).Add(color.FgRed)
	errMessageFmt := color.New(color.FgRed)
	fmt.Println()
	errHeaderFmt.Println("  Error")
	errMessageFmt.Printf("  %s\n", message)
	fmt.Println()
}
