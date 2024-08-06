package connector

import (
	"os/exec"
	"strings"

	"github.com/git-town/git-town/v14/pkg/keys"
)

func loadAccessToken() string {
	process := exec.Command("git", "config", "--get", keys.KeyGithubToken.String()) //nolint:gosec
	output, err := process.Output()
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimSpace(string(output))
}
