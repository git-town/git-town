package script

import (
  "fmt"
  "os"
  "os/exec"
  "strings"

  "github.com/Originate/git-town/lib/git"

  "github.com/fatih/color"
)


func PrintCommand(cmd ...string) {
  header := strings.Join(cmd, " ")
  if strings.HasPrefix(header, "git") {
    header = fmt.Sprintf("[%s] %s", git.GetCurrentBranchName(), header)
  }
  fmt.Println()
  color.New(color.Bold).Println(header)
}


func RunCommand(cmd ...string) error {
  PrintCommand(cmd...)
  subProcess := exec.Command(cmd[0], cmd[1:]...)
  subProcess.Stdout = os.Stdout
  subProcess.Stderr = os.Stderr
  return subProcess.Run()
}
