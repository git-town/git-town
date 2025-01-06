package connector

import (
	"os/exec"
	"strings"

	"github.com/git-town/git-town/v17/pkg"
)

func loadAccessToken() string {
	process := exec.Command("git", "config", "--get", pkg.KeyGithubToken) //nolint:gosec
	output, err := process.Output()
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimSpace(string(output))
}
