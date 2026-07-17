package execute_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/execute"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/giturl"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestNewConnector(t *testing.T) {
	t.Parallel()

	t.Run("no remote URL", func(t *testing.T) {
		t.Parallel()
		repo := execute.OpenRepoResult{} //exhaustruct:ignore
		connector, detectedForgeType, err := repo.NewConnector(None[giturl.Parts]())
		must.NoError(t, err)
		must.Eq(t, None[forgedomain.Connector](), connector)
		must.Eq(t, None[forgedomain.DetectedForgeType](), detectedForgeType)
	})

	t.Run("remote URL at a known forge", func(t *testing.T) {
		t.Parallel()
		repo := execute.OpenRepoResult{} //exhaustruct:ignore
		remoteURL := giturl.Parse("username@github.com:git-town/docs.git")
		connector, detectedForgeType, err := repo.NewConnector(remoteURL)
		must.NoError(t, err)
		must.True(t, connector.IsSome())
		must.Eq(t, Some(forgedomain.ForgeTypeGithub.Detected()), detectedForgeType)
	})

	t.Run("remote URL at an unknown forge", func(t *testing.T) {
		t.Parallel()
		repo := execute.OpenRepoResult{} //exhaustruct:ignore
		remoteURL := giturl.Parse("username@git.example.com:git-town/docs.git")
		connector, detectedForgeType, err := repo.NewConnector(remoteURL)
		must.NoError(t, err)
		must.Eq(t, None[forgedomain.Connector](), connector)
		must.Eq(t, None[forgedomain.DetectedForgeType](), detectedForgeType)
	})
}
