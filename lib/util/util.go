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

// DoesCommandOuputContain runs the given command
// and returns whether its output contains the given string.
func DoesCommandOuputContain(cmd []string, value string) bool {
	return strings.Contains(GetCommandOutput(cmd...), value)
}

// DoesCommandOuputContainLine runs the given command
// and returns whether its output contains teh given string as an entire line.
func DoesCommandOuputContainLine(cmd []string, value string) bool {
	list := strings.Split(GetCommandOutput(cmd...), "\n")
	return DoesStringArrayContain(list, value)
}

// DoesStringArrayContain returns whether the given string slice
// contains the given string.
func DoesStringArrayContain(list []string, value string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

// ExitWithErrorMessage prints the given error message and terminates the application.
func ExitWithErrorMessage(messages ...string) {
	PrintError(messages...)
	os.Exit(1)
}

// GetCommandOutput runs the given command and returns its output.
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

// GetOpenBrowserCommand returns the command to run on the console
// to open the default browser.
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

// GetUserInput reads input from the user and returns it.
func GetUserInput() string {
	text, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatal("Error getting user input:", err)
	}
	return strings.TrimSpace(text)
}

// Indent outputs the given string with the given level of indentation
// on each line. Each level of indentation is two spaces.
func Indent(message string, level int) string {
	prefix := strings.Repeat("  ", level)
	return prefix + strings.Replace(message, "\n", "\n"+prefix, -1)
}

// Pluralize outputs the count and the word. The word is made plural
// if the count isn't one
func Pluralize(count, word string) string {
	result := count + " " + word
	if count != "1" {
		result = result + "s"
	}
	return result
}

// PrintError prints the given error message to the console.
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

// PrintTitle bolds and underlines the given message and prints it to the console
func PrintTitle(title string) {
	titleFmt := color.New(color.Bold).Add(color.Underline)
	titleFmt.Println(title)
}

// RemoveStringFromSlice returns a new string slice which is the given string slice
// with the given string removed
func RemoveStringFromSlice(list []string, value string) (result []string) {
	for _, element := range list {
		if element != value {
			result = append(result, element)
		}
	}
	return
}
