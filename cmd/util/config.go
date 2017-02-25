package util

type Config struct {
  HasRemote bool
  MainBranchName string
  PerennialBranchNames []string
  PullBranchStrategy string
}
