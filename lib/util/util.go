package util

import (
	"bufio"
	"fmt"
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

func ExitWithErrorMessage(messages ...string) {
	PrintError(messages...)
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

var openBrowserCommands = []string{"xdg-open", "open"}
var missingOpenBrowserCommandMessages = []string{
	"Opening a browser requires 'open' on Mac or 'xdg-open' on Linux.",
	"If you would like another command to be supported,",
	"please open an issue at https://github.com/Originate/git-town/issues",
}

func GetOpenBrowserCommand() string {
	for _, command := range openBrowserCommands {
		if GetCommandOutput("which", command) != "" {
			return command
		}
	}
	ExitWithErrorMessage(missingOpenBrowserCommandMessages...)
	return ""
}

var inputReader = bufio.NewReader(os.Stdin)

func GetUserInput() string {
	text, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatal("Error getting user input:", err)
	}
	return strings.TrimSpace(text)
}

func PrintError(messages ...string) {
	errHeaderFmt := color.New(color.Bold).Add(color.FgRed)
	errMessageFmt := color.New(color.FgRed)
	fmt.Println()
	errHeaderFmt.Println("  Error")
	for _, message := range messages {
		errMessageFmt.Println("  " + message)
	}
	fmt.Println()
}
