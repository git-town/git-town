package util

import (
  "fmt"
  "log"
  "os"
  "os/exec"
  "strings"

  "github.com/fatih/color"
)

func DoesStringArrayContain(list []string, value string) bool {
  for _, element := range(list) {
    if element == value {
      return true
    }
  }
  return false
}

func DoesCommandOuputContain(cmd []string, value string) bool {
  return strings.Contains(GetCommandOutput(cmd), value)
}

func DoesCommandOuputContainLine(cmd []string, value string) bool {
  list := strings.Split(GetCommandOutput(cmd), "\n")
  return DoesStringArrayContain(list, value)
}

func ExitWithErrorMessage(message string) {
  errHeaderFmt := color.New(color.Bold).Add(color.FgRed)
  errMessageFmt := color.New(color.FgRed)
  fmt.Println()
  errHeaderFmt.Println("  Error")
  errMessageFmt.Printf("  %s\n", message)
  fmt.Println()
  os.Exit(1)
}

func GetCommandOutput(cmd []string) string {
  subProcess := exec.Command(cmd[0], cmd[1:]...)
  output, err := subProcess.CombinedOutput()
  if err != nil {
    log.Fatal("Command: ", strings.Join(cmd, " "), "\nOutput: " + string(output), "\nError: ", err)
  }
  return strings.TrimSpace(string(output))
}
