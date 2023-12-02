package gitlab_test

import (
	"testing"

	"github.com/git-town/git-town/v10/src/hosting/gitlab"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		config := gitlab.Config{
			Config: common.Config{
				APIToken:     "",
				Hostname:     "gitlab.com",
				Organization: "org",
				Repository:   "repo",
			},
		}
