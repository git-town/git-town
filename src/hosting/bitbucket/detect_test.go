package bitbucket_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/bitbucket"
)

func TestDetect(t *testing.T) {
	t.Parallel()
	t.Run("Bitbucket SaaS", func(t *testing.T) {
		t.Parallel()
		have, err := bitbucket.Detect(configdomain.HostingPlatformNone,
			OriginURL:       giturl.Parse("username@bitbucket.org:git-town/docs.git"),
		})

	})

}
