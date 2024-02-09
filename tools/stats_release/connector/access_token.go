package connector

import (
	"os/exec"
	"strings"
)

func loadAccessToken() string {
	process := exec.Command("git", "config", "--get", "git-town.github-token")
	output, err := process.Output()
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimSpace(string(output))
}
