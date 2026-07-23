package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestIsBitbucket(t *testing.T) {
	t.Parallel()
	tests := map[Option[forgedomain.DetectedForgeType]]bool{
		None[forgedomain.DetectedForgeType]():                     false,
		Some(forgedomain.ForgeTypeBitbucket.Detected()):           true,
		Some(forgedomain.ForgeTypeBitbucketDatacenter.Detected()): true,
		Some(forgedomain.ForgeTypeBitbucketDatacenter.Detected()): true,
		Some(forgedomain.ForgeTypeForgejo.Detected()):             false,
		Some(forgedomain.ForgeTypeGitea.Detected()):               false,
		Some(forgedomain.ForgeTypeGithub.Detected()):              false,
		Some(forgedomain.ForgeTypeGitlab.Detected()):              false,
	}
	for give, want := range tests {
		have := forgedomain.IsBitbucket(give)
		must.Eq(t, want, have)
	}
}
