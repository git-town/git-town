package util

import (
  "fmt"
  "os"
  "os/exec"
  "strings"

  "github.com/fatih/color"
)

func Contains(list []string, item string) bool {
  for i := 0; i < len(list); i++ {
    if item == list[i] {
      return true
    }
  }
  return false
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
  output, _ := subProcess.CombinedOutput()
  return strings.TrimSpace(string(output))
}
