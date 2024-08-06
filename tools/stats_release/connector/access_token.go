package connector

import (
	"os/exec"
	"strings"

	"github.com/git-town/git-town/v14/internal/config/configdomain"
)

func loadAccessToken() string {
	process := exec.Command("git", "config", "--get", configdomain.KeyGithubToken.String()) //nolint:gosec
	output, err := process.Output()
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimSpace(string(output))
}
