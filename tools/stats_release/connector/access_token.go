package connector

import (
	"os/exec"
	"strings"

	"github.com/git-town/git-town/v21/pkg"
)

func loadAccessToken() string {
	process := exec.Command("git", "config", "--get", pkg.KeyGitHubToken) //nolint:gosec
	output, err := process.Output()
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimSpace(string(output))
}
