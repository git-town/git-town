package script

import (
  "fmt"
  "log"
  "os"
  "os/exec"
  "strings"

  "github.com/Originate/gt/cmd/git"

  "github.com/fatih/color"
)

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
    log.Fatal(err)
  }
  return strings.TrimSpace(string(output))
}

func RunCommand(cmd []string) error {
  header := strings.Join(cmd, " ")
  if strings.HasPrefix(header, "git") {
    header = "[" + git.GetCurrentBranchName() + "] " + header
  }
  fmt.Println()
  color.New(color.Bold).Println(header)
  subProcess := exec.Command(cmd[0], cmd[1:]...)
  subProcess.Stdout = os.Stdout
  subProcess.Stderr = os.Stderr
  return subProcess.Run()
}
